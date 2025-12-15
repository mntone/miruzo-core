from pathlib import Path


def validate_relative_path(relpath: Path) -> Path:
	# absolute path is not allowed
	if relpath.is_absolute():
		raise ValueError(f'Absolute path is not allowed: {relpath}')

	# path traversal is not allowed
	if '..' in relpath.parts:
		raise ValueError(f'Path traversal is not allowed: {relpath}')

	# normalize (optional)
	normalized = Path(*relpath.parts)

	return normalized
