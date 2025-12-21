import hashlib
from pathlib import Path


def compute_fingerprint(path: Path) -> str:
	"""Return the SHA-256 hex digest for the file contents."""
	hasher = hashlib.sha256()

	with open(path, 'rb') as file:
		reader = lambda: file.read(4 * 1024 * 1024)

		for chunk in iter(reader, b''):
			hasher.update(chunk)

	fingerprint = hasher.hexdigest()

	return fingerprint
