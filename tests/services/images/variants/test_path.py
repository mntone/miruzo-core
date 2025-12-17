from pathlib import Path

import pytest

from app.services.images.variants.path import (
	_validate_relative_path,
	build_variant_dirpath,
	normalize_relative_path,
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


def test_normalize_relative_path_drops_suffix() -> None:
	result = normalize_relative_path(Path('foo/bar.webp'))
	assert result == Path('foo/bar')


@pytest.mark.parametrize('input_path', [Path('/abs'), Path('..'), Path('foo/../bar')])
def test_normalize_relative_path_raises_for_invalid(input_path: Path) -> None:
	with pytest.raises(ValueError):
		normalize_relative_path(input_path)


def test_build_variant_dirpath_combines_media_root(tmp_path: Path) -> None:
	media_root = tmp_path / 'media'
	(media_root / 'l1w200').mkdir(parents=True)

	result = build_variant_dirpath(media_root, 'l1w200')
	assert result == media_root / 'l1w200'
