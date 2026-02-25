import hashlib
from pathlib import Path

from app.services.ingests.utils.fingerprint import compute_fingerprint


def test_compute_fingerprint_matches_sha256(tmp_path: Path) -> None:
	payload = b'hello world'
	path = tmp_path / 'payload.bin'
	path.write_bytes(payload)

	expected = hashlib.sha256(payload).hexdigest()

	assert compute_fingerprint(path) == expected
