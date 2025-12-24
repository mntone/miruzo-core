from pathlib import Path
from typing import IO

import pytest

from app.utils.files.atomic import ensure_durable_write


def test_ensure_durable_write_writes_content(tmp_path: Path) -> None:
	final_path = tmp_path / 'output.bin'
	payload = b'hello-world'

	def write_fn(handle: IO[bytes]) -> None:
		handle.write(payload)

	ensure_durable_write(final_path, write_fn)

	assert final_path.read_bytes() == payload
	assert list(tmp_path.glob('*.tmp')) == []


def test_ensure_durable_write_removes_temp_on_error(tmp_path: Path) -> None:
	final_path = tmp_path / 'output.bin'

	def write_fn(handle: IO[bytes]) -> None:
		handle.write(b'partial')
		raise RuntimeError('boom')

	with pytest.raises(RuntimeError, match='boom'):
		ensure_durable_write(final_path, write_fn)

	assert not final_path.exists()
	assert list(tmp_path.glob('*.tmp')) == []
