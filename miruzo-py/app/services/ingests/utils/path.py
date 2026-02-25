import unicodedata
from pathlib import Path

from app.config.environments import env

_FORBIDDEN_CHARS = {
	'/',  # POSIX separator
	'\\',  # Windows separator
	':',  # Windows / legacy Mac
	'?',  # glob / Win32
	'<',
	'>',
	'*',
	'|',
	'"',
}


def _validate_origin_relative_path(relpath: Path) -> Path:
	# absolute path
	if relpath.is_absolute():
		raise ValueError(f'Absolute path is not allowed: {relpath}')

	# empty
	if not relpath.parts:
		raise ValueError(f'Empty is not allowed: {relpath}')

	# dot path
	if relpath == Path('.'):
		raise ValueError(f'Dot path is not allowed: {relpath}')

	# path traversal
	if '..' in relpath.parts:
		raise ValueError(f'Path traversal is not allowed: {relpath}')

	for part in relpath.parts:
		# control chars (Unicode Cc)
		if any(unicodedata.category(ch) == 'Cc' for ch in part):
			raise ValueError(f'Control chars are not allowed: {relpath}')

		# forbidden symbols
		if any(ch in part for ch in _FORBIDDEN_CHARS):
			raise ValueError(f'Forbidden character in path: {relpath}')

		# trailing whitespace (Unicode)
		if part and part[-1].isspace():
			raise ValueError(f'Trailing whitespace is not allowed: {relpath}')

		# trailing dot (.)
		if part.endswith('.'):
			raise ValueError(f'Trailing dot is not allowed: {relpath}')

	return relpath


def _ensure_within_root(candidate: Path, *, allowed_root: Path) -> None:
	try:
		candidate.relative_to(allowed_root)
	except ValueError as exc:  # pragma: no cover - defensive
		msg = f'Path {candidate} escapes allowed root {allowed_root}'
		raise ValueError(msg) from exc


def _validate_origin_absolute_path(origin_absolute_path: Path) -> None:
	"""
	Validate that the origin path is located under the configured assets root.

	This ensures the parent directory exists and is within env.gataku_assets_root.
	"""

	resolved_origin_path = origin_absolute_path

	if not resolved_origin_path.is_file():
		raise ValueError(f'origin_path must be a file: {resolved_origin_path}')

	_ensure_within_root(resolved_origin_path, allowed_root=env.gataku_assets_root)


def _map_origin_relative_to_absolute_path(relative_path: Path) -> Path:
	"""
	Build the absolute path for origin assets under the configured root.

	Returns:
		Path: absolute path under assets_root.
	"""

	output_path = env.gataku_assets_root / relative_path

	return output_path


def resolve_origin_absolute_path(relative_path: Path) -> Path:
	_validate_origin_relative_path(relative_path)

	absolute_path = _map_origin_relative_to_absolute_path(relative_path)

	_validate_origin_absolute_path(absolute_path)

	return absolute_path


def map_relative_to_output_path(relative_path: Path) -> Path:
	"""
	Build the output path for copied origin assets.

	Returns:
		Path: output path under media_root.
	"""

	output_path = env.media_root / 'l0orig' / relative_path

	return output_path


def map_relative_to_pathstr(relative_path: Path) -> str:
	"""Build the stored path for copied origin assets (l0orig/... format)."""
	symlink_path = 'l0orig' / relative_path

	return symlink_path.as_posix()


def map_relative_to_symlink_pathstr(relative_path: Path) -> str:
	"""Build the stored path for symlinked origin assets (gataku/... format)."""
	symlink_path = env.gataku_symlink_dirname / relative_path

	return symlink_path.as_posix()
