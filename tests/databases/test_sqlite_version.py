import pytest

from app.databases.sqlite_version import (
	parse_sqlite_version,
	verify_sqlite_supports_returning,
)


@pytest.mark.parametrize(
	('version', 'expected'),
	[
		('3.35.0', (3, 35, 0)),
		('3.45', (3, 45, 0)),
	],
)
def test_parse_sqlite_version(version: str, expected: tuple[int, int, int]) -> None:
	assert parse_sqlite_version(version) == expected


@pytest.mark.parametrize('version', ['3', 'x.35.0', '3.35.beta'])
def test_parse_sqlite_version_raises_for_invalid_format(version: str) -> None:
	with pytest.raises(RuntimeError, match='Unexpected SQLite version format'):
		parse_sqlite_version(version)


def test_verify_sqlite_supports_returning_accepts_boundary() -> None:
	verify_sqlite_supports_returning('3.35.0')


def test_verify_sqlite_supports_returning_rejects_older_version() -> None:
	with pytest.raises(RuntimeError, match='SQLite 3.35.0\\+ is required'):
		verify_sqlite_supports_returning('3.34.1')
