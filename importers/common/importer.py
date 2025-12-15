import json
from datetime import datetime
from logging import getLogger
from pathlib import Path
from shutil import copy2, rmtree
from typing import Literal

from PIL import Image as PILImage
from sqlmodel import Session

from app.core.settings import settings
from app.database import engine, init_database
from app.models.enums import ImageStatus
from app.models.records import ImageRecord, VariantRecord
from app.services.images.thumbnails import (
	VariantReport,
	collect_existing_variants,
	generate_variants,
	reset_variant_directories,
)

log = getLogger(__name__)

ImportMode = Literal['copy', 'symlink']


def ensure_static_root(static_root: Path) -> Path:
	"""Ensure the static root directory exists, creating it if needed."""
	if not static_root.exists():
		static_root.mkdir(parents=True, exist_ok=True)
		print(f'[importer] created static directory: {static_root}')
		return static_root

	if not static_root.is_dir():
		raise RuntimeError(f'Static root must be a directory: {static_root}')

	return static_root


def confirm_overwrite(path: Path, *, force: bool) -> None:
	"""Prompt before deleting populated directories unless force is set."""
	if force:
		return

	choice = input(f'[importer] {path} contains files. Overwrite? [y/N]: ').strip().lower()
	if choice != 'y':
		raise RuntimeError('Aborted by user.')


def prepare_original_dir(
	static_root: Path,
	mode: ImportMode,
	original_subdir: str,
	assets_root: Path,
	*,
	force: bool,
) -> Path:
	"""Set up static/orig to either symlink to the assets root or act as a real directory."""
	original_dir = static_root / original_subdir

	if original_dir.exists():
		if original_dir.is_symlink():
			current_target = original_dir.resolve(strict=False)
			print(f'[importer] {original_dir} symlink detected -> {current_target}')
			if mode == 'symlink':
				if current_target == assets_root:
					return original_dir
				if not force:
					choice = (
						input(
							f'[importer] {original_dir} points to {current_target}, expected {assets_root}. Recreate? [y/N]: ',
						)
						.strip()
						.lower()
					)
					if choice != 'y':
						raise RuntimeError('Aborted due to mismatched static/orig symlink.')
				original_dir.unlink()
			else:
				if not force:
					choice = (
						input(
							f'[importer] {original_dir} is a symlink but copy mode was requested. Recreate directory? [y/N]: ',
						)
						.strip()
						.lower()
					)
					if choice != 'y':
						raise RuntimeError('Aborted due to symlink/copy mode mismatch.')
				original_dir.unlink()
		elif original_dir.is_dir():
			if any(original_dir.iterdir()):
				confirm_overwrite(original_dir, force=force)
			rmtree(original_dir)
		else:
			confirm_overwrite(original_dir, force=force)
			original_dir.unlink()

	if mode == 'symlink':
		original_dir.symlink_to(assets_root, target_is_directory=True)
		print(f'[importer] linked {original_dir} -> {assets_root}')
	else:
		original_dir.mkdir(parents=True, exist_ok=True)
		print(f'[importer] prepared copy directory: {original_dir}')

	return original_dir


def import_jsonl(
	jsonl_path: str,
	static_dir: str,
	limit: int = 100,
	mode: ImportMode = 'symlink',
	original_subdir: str = 'gataku',
	force: bool = False,
	report_variants: bool = False,
	repair: bool = False,
) -> None:
	"""Read gataku JSONL data, populate the database, and copy/symlink assets plus thumbnails."""
	static_root = ensure_static_root(Path(static_dir))
	gataku_root = settings.gataku_root.resolve()
	assets_root = settings.assets_root.resolve()
	orig_dir = prepare_original_dir(static_root, mode, original_subdir, assets_root, force=force)

	variant_layers = settings.variant_layers

	if repair:
		print('[importer] repair mode enabled: skipping thumbnail generation.')
	else:
		reset_variant_directories(static_root, variant_layers)
	init_database()

	session = Session(engine)

	stats = {
		'read': 0,
		'imported': 0,
		'invalid': 0,
		'missing': 0,
	}

	with open(jsonl_path, 'r') as f:
		for i, line in enumerate(f):
			if i >= limit:
				break

			stats['read'] += 1

			try:
				record = json.loads(line)
			except Exception:
				stats['invalid'] += 1
				continue  # skip invalid JSON

			raw_path = Path(record['filepath'])
			src_path = raw_path if raw_path.is_absolute() else gataku_root / raw_path
			if not src_path.exists():
				log.warning(f'missing file: {src_path}')
				stats['missing'] += 1
				continue  # image not found, skip entry

			try:
				original_size = src_path.stat().st_size
			except OSError:
				original_size = None

			try:
				relative_asset_path = src_path.relative_to(assets_root)
			except ValueError:
				log.warning(f'file outside assets root: {src_path}')
				stats['missing'] += 1
				continue

			if mode == 'copy':
				dst = orig_dir / relative_asset_path
				if not dst.exists():
					dst.parent.mkdir(parents=True, exist_ok=True)
					copy2(src_path, dst)

			width = None
			height = None
			original_variant = None
			variant_records: list[list[VariantRecord]] = []
			variant_reports: list[VariantReport] = []
			try:
				with PILImage.open(src_path) as pil_image:
					width = pil_image.width
					height = pil_image.height
					mime_type = PILImage.MIME.get(pil_image.format, 'application/octet-stream')
					format_name, codecs = _map_original_variant_attrs(pil_image, mime_type)
					public_path = f'/static/{original_subdir}/{relative_asset_path.as_posix()}'
					original_variant = {
						'filepath': public_path,
						'format': format_name,
						'codecs': codecs,
						'size': original_size,
						'width': width,
						'height': height,
						'quality': None,
					}
					if not repair:
						variant_records, variant_reports = generate_variants(
							pil_image,
							relative_asset_path,
							static_root,
							variant_layers,
							original_size=original_size,
						)
			except Exception as exc:
				log.warning('failed to prepare metadata/thumbnail for %s: %s', src_path, exc)
				if repair:
					variant_records = collect_existing_variants(
						relative_asset_path,
						static_root,
						variant_layers,
					)

			if repair and not variant_records:
				variant_records = collect_existing_variants(
					relative_asset_path,
					static_root,
					variant_layers,
				)
			if repair:
				variant_reports = []

			if report_variants and (variant_reports or original_size):
				print_variant_report(relative_asset_path, original_size, variant_reports)

			if original_size is None or width is None or height is None:
				log.warning('skipping %s due to missing size/width/height metadata', src_path)
				continue

			if original_variant is None:
				public_path = f'/static/{original_subdir}/{relative_asset_path.as_posix()}'
				original_variant = {
					'filepath': public_path,
					'format': 'unknown',
					'codecs': None,
					'size': original_size,
					'width': width,
					'height': height,
					'quality': None,
				}

			captured_at = None
			created_at_value = record.get('created_at')
			if created_at_value:
				try:
					captured_at = datetime.fromisoformat(created_at_value)
				except Exception:
					pass

			img = ImageRecord(
				fingerprint=record['sha256'],
				captured_at=captured_at,
				status=ImageStatus.ACTIVE,
				original=original_variant,
				fallback=None,
				variants=variant_records,
			)

			session.add(img)
			stats['imported'] += 1

			if stats['read'] % 10 == 0 or stats['read'] == limit:
				print(
					f'[importer] progress: read={stats["read"]}, imported={stats["imported"]}, invalid={stats["invalid"]}, missing={stats["missing"]}',
				)
		session.commit()
	session.close()

	print(
		f'[importer] summary: read={stats["read"]}, imported={stats["imported"]}, invalid={stats["invalid"]}, missing={stats["missing"]}',
	)


def _format_bytes(size: int) -> str:
	"""Convert a size in bytes into a human-friendly string."""
	if size is None:
		return 'n/a'

	thresholds = [
		(1024**3, 'GB'),
		(1024**2, 'MB'),
		(1024, 'KB'),
	]
	for factor, suffix in thresholds:
		if size >= factor:
			value = size / factor
			if value < 10:
				return f'{value:.2f} {suffix}'
			return f'{value:.1f} {suffix}'
	return f'{size} B'


def print_variant_report(
	asset_path: Path,
	original_size: int | None,
	reports: list[VariantReport],
) -> None:
	print(f'[importer] variant report ({asset_path.as_posix()}):')
	header = f'{"Layer":<8} {"Label":<8} {"Resolution":<12} {"Size":>12} {"Ratio":<12}'
	print(header)
	print('-' * len(header))

	if original_size:
		print(
			f'{"original":<8} {"orig":<8} {"-":<12} {_format_bytes(original_size):>12} {"n/a":<12}',
		)

	for report in reports:
		size_str = _format_bytes(report.size_bytes or 0)
		if (
			report.ratio_percent is not None
			and report.delta_bytes is not None
			and original_size
			and original_size > 0
		):
			diff_str = _format_bytes(abs(report.delta_bytes))
			sign = '+' if report.delta_bytes >= 0 else '-'
			ratio_str = f'{report.ratio_percent:.0f}% ({sign}{diff_str})'
		else:
			ratio_str = 'n/a'

		resolution = f'{report.width}x{report.height}'
		print(
			f'{report.layer_name:<8} {report.label:<8} {resolution:<12} {size_str:>12} {ratio_str:<12}',
		)


def _map_original_variant_attrs(pil_image: PILImage.Image, mime_type: str) -> tuple[str, str | None]:
	fmt = (pil_image.format or '').upper()
	if fmt == 'BMP':
		return 'bmp', None
	if fmt == 'GIF':
		return 'gif', None
	if fmt == 'PNG':
		return 'png', None
	if fmt in {'JPG', 'JPEG'}:
		return 'jpeg', None
	if fmt == 'TIFF':
		return 'tiff', None
	if fmt == 'WEBP':
		if pil_image.info.get('lossless'):
			return 'webp', 'vp8l'
		return 'webp', 'vp8'
	return mime_type.split('/')[-1], None
