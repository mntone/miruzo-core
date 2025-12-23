import unicodedata
from pathlib import Path
from typing import NewType

from app.config.variant import VariantSlotkey

OriginRelativePath = NewType('OriginRelativePath', Path)
VariantRelativePath = NewType('VariantRelativePath', Path)

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


def _validate_relative_path(relpath: Path) -> Path:
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


def build_origin_relative_path(relative_path: Path) -> OriginRelativePath:
	relpath_noext = relative_path.with_suffix('')

	validated_relpath = _validate_relative_path(relpath_noext)

	return OriginRelativePath(validated_relpath)


def build_variant_relative_path(
	origin_path: OriginRelativePath,
	*,
	under: str | VariantSlotkey,
) -> VariantRelativePath:
	if isinstance(under, VariantSlotkey):
		variant_dirname = under.label
	else:
		variant_dirname = under

	media_relpath = Path(variant_dirname).joinpath(origin_path)
	return VariantRelativePath(media_relpath)


def build_absolute_path(variant_relpath: VariantRelativePath, *, under: Path) -> Path:
	# Normalize argument name for internal use
	media_root = under

	absolute_path = media_root.joinpath(variant_relpath)
	return absolute_path
