from pathlib import Path

import pytest

from app.services.images.variants.path import VariantRelativePath
from app.services.images.variants.types import FileInfo


def test_file_info_from_relative_path_reads_stat(tmp_path: Path) -> None:
	relative_path = VariantRelativePath(Path('l1w320/foo/bar.webp'))
	absolute_path = tmp_path / relative_path
	absolute_path.parent.mkdir(parents=True, exist_ok=True)
	absolute_path.write_bytes(b'data')

	info = FileInfo.from_relative_path(relative_path, under=tmp_path)

	assert info.absolute_path == absolute_path
	assert info.relative_path == relative_path
	assert info.bytes == 4


def test_file_info_from_relative_path_raises_when_missing(tmp_path: Path) -> None:
	relative_path = VariantRelativePath(Path('l1w320/missing.webp'))

	with pytest.raises(FileNotFoundError):
		FileInfo.from_relative_path(relative_path, under=tmp_path)
