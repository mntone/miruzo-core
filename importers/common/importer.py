from pathlib import Path
from shutil import rmtree

from sqlmodel import Session

from importers.common.ingest_time import resolve_captured_at
from importers.common.origin import OriginResolver
from importers.common.readers.jsonl import JsonlReader
from importers.common.report import ImportStats, ProgressReporter

from app.config.environments import Settings
from app.config.environments import env as global_env
from app.databases import engine, init_database
from app.models.enums import IngestMode
from app.persist.images.factory import create_image_repository
from app.persist.ingests.factory import create_ingest_repository
from app.services.images.ingest import ImageIngestService
from app.services.images.variants.bootstrap import configure_pillow
from app.services.images.variants.types import DEFAULT_VARIANT_POLICY
from app.services.ingests.bootstrap import ensure_ingest_layout


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

	gataku_assets_root = env.gataku_assets_root

	configure_pillow()
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

	stats = ImportStats()
	warned_created_at_fallback = False
	reporter = ProgressReporter(report_variants=report_variants)
	reader = JsonlReader(Path(jsonl_path))
	resolver = OriginResolver(
		gataku_root=env.gataku_root,
		gataku_assets_root=gataku_assets_root,
	)

	with session:
		for row in reader.read(limit=limit):
			stats.read += 1

			if row is None:
				stats.invalid += 1
				continue  # skip invalid JSON

			resolution = resolver.resolve(row)
			if resolution is None:
				stats.missing += 1
				continue  # image not found or invalid path

			src_path = resolution.src_path
			origin_relative_path = resolution.origin_relative_path

			captured_at, used_fallback, warned_created_at_fallback = resolve_captured_at(
				created_at_value=row.created_at,
				src_path=src_path,
				warned_fallback=warned_created_at_fallback,
			)
			if used_fallback:
				stats.fallback += 1

			ingest_record = ingest.ingest(
				origin_path=origin_relative_path,
				fingerprint=row.sha256,
				captured_at=captured_at,
				ingest_mode=mode,
			)
			stats.ingested += 1

			reporter.maybe_report_variants(ingest_record)
			reporter.report_progress(stats, force=stats.read == limit)

	reporter.report_summary(stats)
