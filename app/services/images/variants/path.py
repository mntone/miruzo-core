from pathlib import Path
from typing import NewType

from app.config.variant import VariantSlotkey

VariantBasePath = NewType('VariantBasePath', Path)
VariantRelativePath = NewType('VariantRelativePath', Path)


def map_origin_to_variant_basepath(relative_path: Path) -> VariantBasePath:
	"""Drop the origin prefix and suffix to build the variant base path."""
	parts = relative_path.parts
	if len(parts) < 2:
		raise ValueError(f'Origin path must include a prefix: {relative_path}')

	relpath_noext = Path(*parts[1:]).with_suffix('')

	return VariantBasePath(relpath_noext)


def build_variant_relative_path(
	variant_basepath: VariantBasePath,
	*,
	under: str | VariantSlotkey,
) -> VariantRelativePath:
	if isinstance(under, VariantSlotkey):
		variant_dirname = under.label
	else:
		variant_dirname = under

	media_relpath = Path(variant_dirname).joinpath(variant_basepath)
	return VariantRelativePath(media_relpath)


def build_absolute_path(variant_relpath: VariantRelativePath, *, under: Path) -> Path:
	# Normalize argument name for internal use
	media_root = under

	absolute_path = media_root.joinpath(variant_relpath)
	return absolute_path
