import shutil
from pathlib import Path


def copy_origin_file(
	src: Path,
	dst: Path,
) -> None:
	"""
	Copy the original file to the ingest destination, creating parents as needed.

	Raises:
		FileNotFoundError: when either the original file or destination parent is missing.
		IOError/OSError: propagated from shutil.copy2 on failure.
	"""

	if not src.is_file():
		raise ValueError(f'Origin path is not a file: {src}')

	parent = dst.parent
	parent.mkdir(parents=True, exist_ok=True)

	shutil.copy2(src, dst)


def delete_origin_file(path: Path) -> None:
	"""
	Delete an origin file created during ingest.
	"""

	if not path.exists():
		return
	if not path.is_file():
		raise ValueError(f'Origin path is not a file: {path}')

	path.unlink()
