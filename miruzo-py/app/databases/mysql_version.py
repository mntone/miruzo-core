import re

_MYSQL_VERSION_PATTERN = re.compile(r'^(\d+)\.(\d+)(?:\.(\d+))?')


def parse_mysql_version(version: str) -> tuple[int, int, int]:
	match = _MYSQL_VERSION_PATTERN.match(version)
	if match is None:
		raise RuntimeError(f'Unexpected MySQL version format: {version}')

	major = int(match.group(1))
	minor = int(match.group(2))
	patch = int(match.group(3)) if match.group(3) is not None else 0
	return major, minor, patch


def verify_mysql_supports_check_constraints(version: str) -> None:
	parsed = parse_mysql_version(version)
	min_required = (8, 0, 16)
	if parsed < min_required:
		raise RuntimeError(
			f'MySQL 8.0.16+ is required for CHECK constraint support: detected {version}',
		)
