from pathlib import Path

import pytest

from app.services.images.variants.security import validate_relative_path


def test_validate_relative_path_accepts_simple_relative_paths() -> None:
	path = validate_relative_path(Path('foo/bar/baz'))
	assert path == Path('foo/bar/baz')


@pytest.mark.parametrize('input_path', [Path('/etc/passwd'), Path('/foo/bar')])
def test_validate_relative_path_rejects_absolute(input_path: Path) -> None:
	with pytest.raises(ValueError):
		validate_relative_path(input_path)


@pytest.mark.parametrize('input_path', [Path('../foo'), Path('foo/../bar'), Path('..')])
def test_validate_relative_path_rejects_path_traversal(input_path: Path) -> None:
	with pytest.raises(ValueError):
		validate_relative_path(input_path)
