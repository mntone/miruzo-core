from pathlib import Path

from app.config.variant import VariantSlot
from app.services.images.variants.path import (
	build_absolute_path,
	build_variant_relative_path,
	map_origin_to_variant_basepath,
)


def test_map_origin_to_variant_basepath_drops_prefix_and_suffix() -> None:
	result = map_origin_to_variant_basepath(Path('l0orig/foo/bar.webp'))
	assert result == Path('foo/bar')


def test_build_variant_relative_path_uses_slot() -> None:
	basepath = map_origin_to_variant_basepath(Path('l0orig/foo/bar.webp'))
	slot = VariantSlot(layer_id=1, width=200)

	relpath = build_variant_relative_path(basepath, under=slot)

	assert relpath == Path('l1w200/foo/bar')


def test_build_variant_relative_path_uses_string() -> None:
	basepath = map_origin_to_variant_basepath(Path('l0orig/foo/bar.webp'))
	variant_dirname = 'l1w200'

	relpath = build_variant_relative_path(basepath, under=variant_dirname)

	assert relpath == Path('l1w200/foo/bar')


def test_build_absolute_path_combines_media_root(tmp_path: Path) -> None:
	basepath = map_origin_to_variant_basepath(Path('l0orig/foo/bar.webp'))
	relpath = build_variant_relative_path(basepath, under='l1w200')

	result = build_absolute_path(relpath, under=tmp_path)

	assert result == tmp_path / 'l1w200/foo/bar'
