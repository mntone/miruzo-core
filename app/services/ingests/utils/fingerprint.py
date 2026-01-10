import hashlib
import re
from pathlib import Path

_HEX_DIGEST_RE = re.compile(r'^[0-9a-fA-F]{64}$')


def compute_fingerprint(path: Path) -> str:
	"""Return the SHA-256 hex digest for the file contents."""
	hasher = hashlib.sha256()

	with open(path, 'rb') as file:
		reader = lambda: file.read(4 * 1024 * 1024)

		for chunk in iter(reader, b''):
			hasher.update(chunk)

	fingerprint = hasher.hexdigest()

	return fingerprint


def normalize_fingerprint(value: str) -> str | None:
	"""Return a normalized SHA-256 hex digest, or None when invalid."""

	normalized = value.strip().lower()
	if not _HEX_DIGEST_RE.fullmatch(normalized):
		return None

	return normalized
