import os
import sys
from pathlib import Path
from tempfile import NamedTemporaryFile
from typing import IO, Callable

if sys.platform == 'win32':
	import ctypes

	FILE_ATTRIBUTE_HIDDEN = 0x02
	FILE_ATTRIBUTE_NORMAL = 0x80

	def _set_hidden_if_required(path: Path) -> None:
		ctypes.windll.kernel32.SetFileAttributesW(
			path.__str__(),
			FILE_ATTRIBUTE_HIDDEN,
		)

	def _unset_hidden_if_required(path: Path) -> None:
		ctypes.windll.kernel32.SetFileAttributesW(
			path.__str__(),
			FILE_ATTRIBUTE_NORMAL,
		)

	def _fsync_dir(dir_path: Path) -> None:
		pass

	def _get_file_prefix(path: Path) -> str:
		return path.name
else:

	def _set_hidden_if_required(path: Path) -> None:
		pass

	def _unset_hidden_if_required(path: Path) -> None:
		pass

	def _fsync_dir(dir_path: Path) -> None:
		fd = os.open(dir_path, os.O_DIRECTORY)
		try:
			os.fsync(fd)
		finally:
			os.close(fd)

	def _get_file_prefix(path: Path) -> str:
		return '.' + path.name


def ensure_durable_write(
	final_path: Path,
	write_fn: Callable[[IO[bytes]], None],
) -> None:
	"""Atomically write a file and fsync it; parent directory must exist."""

	tmp_path: Path | None = None
	try:
		# same directory, hidden tmp
		with NamedTemporaryFile(
			mode='wb',
			suffix='.tmp',
			prefix=_get_file_prefix(final_path),
			dir=final_path.parent,
			delete=False,
		) as tmp:
			tmp_path = Path(tmp.name)

			# set hidden attr on Windows
			_set_hidden_if_required(tmp_path)

			# write via file descriptor
			write_fn(tmp)

			# flush python buffer
			tmp.flush()

			# fsync file
			os.fsync(tmp.fileno())

		# atomic rename
		os.replace(tmp_path, final_path)

		# unset hidden attr on Windows
		_unset_hidden_if_required(final_path)

	except Exception:
		if tmp_path is not None:
			tmp_path.unlink(missing_ok=True)
		raise

	# fsync directory
	_fsync_dir(final_path.parent)
