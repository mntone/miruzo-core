from pathlib import Path

from app.config.environments import env


def _ensure_within_root(candidate: Path, *, allowed_root: Path) -> None:
	try:
		candidate.relative_to(allowed_root)
	except ValueError as exc:  # pragma: no cover - defensive
		msg = f'Path {candidate} escapes allowed root {allowed_root}'
		raise ValueError(msg) from exc


def validate_origin_path(origin_path: Path) -> Path:
	"""
	Validate that the origin path is located under the configured assets root.

	This ensures the parent directory exists and is within env.gataku_assets_root.
	"""

	resolved_origin_path = origin_path.resolve(strict=True)

	if not resolved_origin_path.is_file():
		raise ValueError(f'origin_path must be a file: {resolved_origin_path}')

	_ensure_within_root(resolved_origin_path, allowed_root=env.gataku_assets_root)

	return resolved_origin_path


def map_origin_to_paths(origin_path: Path) -> tuple[str, Path]:
	"""
	Build the media-relative path and output path for copied origin assets.

	Returns:
		tuple[str, Path]: media-relative path (l0orig/...) and output path under media_root.
	"""

	origin_root = env.gataku_assets_root.resolve()
	relative_origin_path = origin_path.relative_to(origin_root)

	relative_path = 'l0orig/' + relative_origin_path.as_posix()
	output_path = env.media_root / 'l0orig' / relative_origin_path

	return relative_path, output_path


def map_origin_to_symlink_paths(origin_path: Path) -> tuple[str, Path]:
	"""
	Build the media-relative path and output path for symlinked origin assets.

	Returns:
		tuple[str, Path]: media-relative path (gataku/...) and output path under media_root.
	"""

	origin_root = env.gataku_assets_root.resolve()
	relative_origin_path = origin_path.relative_to(origin_root)

	relative_path = f'{env.gataku_symlink_dirname}/{relative_origin_path.as_posix()}'
	output_path = env.media_root / env.gataku_symlink_dirname / relative_origin_path

	return relative_path, output_path
