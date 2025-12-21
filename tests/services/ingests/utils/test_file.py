from pathlib import Path

import pytest

from app.services.ingests.utils.file import copy_origin_file, delete_origin_file


def test_copy_origin_file_copies_bytes(tmp_path: Path) -> None:
	src = tmp_path / 'src.txt'
	src.write_bytes(b'hello')
	dst = tmp_path / 'nested' / 'dst.txt'

	copy_origin_file(src, dst)

	assert dst.exists()
	assert dst.read_bytes() == b'hello'


def test_copy_origin_file_rejects_non_file(tmp_path: Path) -> None:
	src = tmp_path / 'dir'
	src.mkdir()
	dst = tmp_path / 'dst.txt'

	with pytest.raises(ValueError, match='not a file'):
		copy_origin_file(src, dst)


def test_delete_origin_file_removes_existing(tmp_path: Path) -> None:
	target = tmp_path / 'delete.txt'
	target.write_text('data')

	delete_origin_file(target)

	assert not target.exists()


def test_delete_origin_file_ignores_missing(tmp_path: Path) -> None:
	target = tmp_path / 'missing.txt'

	delete_origin_file(target)

	assert not target.exists()


def test_delete_origin_file_rejects_non_file(tmp_path: Path) -> None:
	target = tmp_path / 'dir'
	target.mkdir()

	with pytest.raises(ValueError, match='not a file'):
		delete_origin_file(target)
