from pathlib import Path

import pytest

from app.config.variant import VariantSlotkey
from app.services.images.variants.path import (
	_validate_relative_path,
	build_absolute_path,
	build_origin_relative_path,
	build_variant_relative_path,
)


def test_validate_relative_path_accepts_simple_relative_paths() -> None:
	path = _validate_relative_path(Path('foo/bar/baz'))
	assert path == Path('foo/bar/baz')


@pytest.mark.parametrize('input_path', [Path('/etc/passwd'), Path('/foo/bar')])
def test_validate_relative_path_rejects_absolute(input_path: Path) -> None:
	with pytest.raises(ValueError):
		_validate_relative_path(input_path)


@pytest.mark.parametrize('input_path', [Path('../foo'), Path('foo/../bar'), Path('..')])
def test_validate_relative_path_rejects_path_traversal(input_path: Path) -> None:
	with pytest.raises(ValueError):
		_validate_relative_path(input_path)


@pytest.mark.parametrize(
	'input_path',
	[
		Path(''),
		Path('.'),
		Path('foo\x00bar'),
		Path('foo|bar'),
		Path('foo '),
		Path('foo\u3000'),
		Path('foo.'),
	],
)
def test_validate_relative_path_rejects_empty_dot_and_invalid_chars(input_path: Path) -> None:
	with pytest.raises(ValueError):
		_validate_relative_path(input_path)


def test_build_origin_relative_path_drops_suffix() -> None:
	result = build_origin_relative_path(Path('foo/bar.webp'))
	assert result == Path('foo/bar')


@pytest.mark.parametrize('input_path', [Path('/abs'), Path('..'), Path('foo/../bar')])
def test_build_origin_relative_path_raises_for_invalid(input_path: Path) -> None:
	with pytest.raises(ValueError):
		build_origin_relative_path(input_path)


def test_build_variant_relative_path_uses_slotkey() -> None:
	origin = build_origin_relative_path(Path('foo/bar.webp'))
	slotkey = VariantSlotkey(layer_id=1, width=200)

	relpath = build_variant_relative_path(origin, under=slotkey)

	assert relpath == Path('l1w200/foo/bar')


def test_build_variant_relative_path_uses_string() -> None:
	origin = build_origin_relative_path(Path('foo/bar.webp'))
	variant_dirname = 'l1w200'

	relpath = build_variant_relative_path(origin, under=variant_dirname)

	assert relpath == Path('l1w200/foo/bar')


def test_build_absolute_path_combines_media_root(tmp_path: Path) -> None:
	origin = build_origin_relative_path(Path('foo/bar.webp'))
	relpath = build_variant_relative_path(origin, under='l1w200')

	result = build_absolute_path(relpath, under=tmp_path)

	assert result == tmp_path / 'l1w200/foo/bar'
