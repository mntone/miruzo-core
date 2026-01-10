import json
from logging import getLogger
from pathlib import Path
from shutil import rmtree

from sqlmodel import Session

from importers.common.ingest_time import resolve_captured_at

from app.config.environments import Settings
from app.config.environments import env as global_env
from app.database import engine, init_database
from app.models.enums import IngestMode
from app.models.records import IngestRecord
from app.persist.images.factory import create_image_repository
from app.persist.ingests.factory import create_ingest_repository
from app.services.images.ingest import ImageIngestService
from app.services.images.variants.types import DEFAULT_VARIANT_POLICY
from app.services.ingests.bootstrap import ensure_ingest_layout

log = getLogger(__name__)


def confirm_overwrite(path: Path, *, force: bool) -> None:
	"""Prompt before deleting populated directories unless force is set."""
	if force:
		return

	choice = input(f'[importer] {path} contains files. Overwrite? [y/N]: ').strip().lower()
	if choice != 'y':
		raise RuntimeError('Aborted by user.')


def prepare_original_directory(
	*,
	media_root: Path,
	gataku_assets_root: Path,
	gataku_symlink_dirname: str,
	mode: IngestMode,
	force: bool,
) -> None:
	"""Set up media/orig to either symlink to the assets root or act as a real directory."""

	if mode != IngestMode.SYMLINK:
		return

	original_dir = media_root / gataku_symlink_dirname
	if original_dir.is_symlink() and not original_dir.exists():
		raise RuntimeError(f'Dangling symlink detected: {original_dir}')
	elif original_dir.is_symlink():
		current_target = original_dir.resolve()
		print(f'[importer] {original_dir} symlink detected -> {current_target}')
		if current_target == gataku_assets_root:
			return

		if not force:
			choice = (
				input(
					f'[importer] {original_dir} points to {current_target}, expected {gataku_assets_root}. Recreate? [y/N]: ',
				)
				.strip()
				.lower()
			)
			if choice != 'y':
				raise RuntimeError('Aborted due to mismatched media/orig symlink.')
		original_dir.unlink()
	elif original_dir.exists():
		if original_dir.is_dir():
			if any(original_dir.iterdir()):
				confirm_overwrite(original_dir, force=force)
			rmtree(original_dir)
		else:
			confirm_overwrite(original_dir, force=force)
			original_dir.unlink()

	original_dir.symlink_to(gataku_assets_root, target_is_directory=True)
	print(f'[importer] linked {original_dir} -> {gataku_assets_root}')


def import_jsonl(
	jsonl_path: str,
	limit: int,
	mode: IngestMode,
	force: bool,
	report_variants: bool = False,
	env: Settings = global_env,
) -> None:
	"""Read gataku JSONL data, populate the database, and copy/symlink assets plus thumbnails."""

	gataku_root = env.gataku_root
	gataku_assets_root = env.gataku_assets_root

	ensure_ingest_layout(env)
	prepare_original_directory(
		media_root=env.media_root,
		gataku_assets_root=gataku_assets_root,
		gataku_symlink_dirname=env.gataku_symlink_dirname,
		mode=mode,
		force=force,
	)
	init_database()

	session = Session(engine)
	image_repo = create_image_repository(session)
	ingest_repo = create_ingest_repository(session)
	ingest = ImageIngestService(
		image_repo=image_repo,
		ingest_repo=ingest_repo,
		policy=DEFAULT_VARIANT_POLICY,
	)

	stats = {
		'read': 0,
		'ingested': 0,
		'invalid': 0,
		'missing': 0,
		'fallback': 0,
	}
	warned_created_at_fallback = False

	with session:
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
					origin_relative_path = src_path.relative_to(gataku_assets_root)
				except ValueError:
					log.warning(f'file outside assets root: {src_path}')
					stats['missing'] += 1
					continue

				captured_at, used_fallback, warned_created_at_fallback = resolve_captured_at(
					created_at_value=record.get('created_at'),
					src_path=src_path,
					warned_fallback=warned_created_at_fallback,
				)
				if used_fallback:
					stats['fallback'] += 1

				ingest_record = ingest.ingest(
					origin_path=origin_relative_path,
					fingerprint=record['sha256'],
					captured_at=captured_at,
					ingest_mode=mode,
				)
				stats['ingested'] += 1

				if report_variants:
					print_variant_report(ingest_record)

				if stats['read'] % 10 == 0 or stats['read'] == limit:
					print(
						f'[importer] progress: read={stats["read"]}, ingested={stats["ingested"]}, invalid={stats["invalid"]}, missing={stats["missing"]}, fallback={stats["fallback"]}',
					)

	print(
		f'[importer] summary: read={stats["read"]}, ingested={stats["ingested"]}, invalid={stats["invalid"]}, missing={stats["missing"]}, fallback={stats["fallback"]}',
	)


def _format_bytes(size: int) -> str:
	"""Convert a size in bytes into a human-friendly string."""

	thresholds = [
		(1024**4, 'TB'),
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


def print_variant_report(record: IngestRecord) -> None:
	image = record.image
	if image is None:
		return

	print(f'[importer] variant report ({record.relative_path}):')
	header = f'{"Label":<10} {"Resolution":<12} {"Size":>10} {"Ratio":<12}'
	print(header)
	print('-' * len(header))

	original = image.original
	original_width = original['width']
	original_height = original['height']
	original_resolution = f'{original_width}x{original_height}'
	original_size = original['bytes']
	print(
		f'{"l0orig":<10} {original_resolution:<12} {_format_bytes(original_size):>10} {"n/a":<12}',
	)

	for variant in image.variants:
		label = 'l' + variant['layer_id'].__str__() + 'w' + variant['width'].__str__()
		width = variant['width']
		height = variant['height']

		size = variant['bytes']
		size_str = _format_bytes(size or 0)

		delta = size - original_size
		delta_str = _format_bytes(abs(delta))

		ratio = (size / original_size) * 100
		ratio_sign = '+' if delta >= 0 else '-'
		ratio_str = f'{ratio:.1f}% ({ratio_sign}{delta_str})'

		resolution = f'{width}x{height}'
		print(f'{label:<10} {resolution:<12} {size_str:>10} {ratio_str:<12}')
