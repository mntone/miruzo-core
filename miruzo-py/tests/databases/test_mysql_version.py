import pytest

from app.databases.mysql_version import (
	parse_mysql_version,
	verify_mysql_supports_check_constraints,
)


@pytest.mark.parametrize(
	('version', 'expected'),
	[
		('8.0.16', (8, 0, 16)),
		('8.4', (8, 4, 0)),
		('8.4.6-0ubuntu0.24.04.1', (8, 4, 6)),
		('9.0.1 MySQL Community Server - GPL', (9, 0, 1)),
	],
)
def test_parse_mysql_version(version: str, expected: tuple[int, int, int]) -> None:
	assert parse_mysql_version(version) == expected


@pytest.mark.parametrize('version', ['8', 'x.0.16', 'mysql-8.0.16'])
def test_parse_mysql_version_raises_for_invalid_format(version: str) -> None:
	with pytest.raises(RuntimeError, match='Unexpected MySQL version format'):
		parse_mysql_version(version)


def test_verify_mysql_supports_check_constraints_accepts_boundary() -> None:
	verify_mysql_supports_check_constraints('8.0.16')


def test_verify_mysql_supports_check_constraints_rejects_older_version() -> None:
	with pytest.raises(RuntimeError, match='MySQL 8.0.16\\+ is required'):
		verify_mysql_supports_check_constraints('8.0.15')
