import stat
from pathlib import Path


def ensure_directory_access(path: Path, name: str | None = None) -> None:
	name = name or 'path'

	try:
		st = path.stat()
	except FileNotFoundError as exc:
		raise RuntimeError(name + ' does not exist') from exc

	if not stat.S_ISDIR(st.st_mode):
		raise RuntimeError(name + ' is not a directory')
