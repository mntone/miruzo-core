def parse_sqlite_version(version: str) -> tuple[int, int, int]:
	parts = version.split('.')
	if len(parts) < 2:
		raise RuntimeError(f'Unexpected SQLite version format: {version}')

	try:
		major = int(parts[0])
		minor = int(parts[1])
		patch = int(parts[2]) if len(parts) >= 3 else 0
	except ValueError as exc:
		raise RuntimeError(f'Unexpected SQLite version format: {version}') from exc

	return major, minor, patch


def verify_sqlite_supports_returning(version: str) -> None:
	parsed = parse_sqlite_version(version)
	min_required = (3, 35, 0)
	if parsed < min_required:
		raise RuntimeError(
			f'SQLite 3.35.0+ is required for RETURNING support: detected {version}',
		)
