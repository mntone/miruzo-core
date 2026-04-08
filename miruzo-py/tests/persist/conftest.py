from collections.abc import Iterator

import pytest
from sqlalchemy.orm import Session

from tests.persist.fixtures.backends.mysql import mysql_dsn, mysql_session
from tests.persist.fixtures.backends.postgres import postgres_dsn, postgres_session
from tests.persist.fixtures.backends.sqlite import sqlite_session

from app.config.environments import DatabaseBackend

__all__ = [
	'mysql_dsn',
	'mysql_session',
	'postgres_dsn',
	'postgres_session',
	'sqlite_session',
]


@pytest.fixture()
def session(request: pytest.FixtureRequest) -> Iterator[Session]:
	backend = getattr(request, 'param', DatabaseBackend.SQLITE)
	match backend:
		case DatabaseBackend.MYSQL:
			yield request.getfixturevalue('mysql_session')
		case DatabaseBackend.POSTGRE_SQL:
			yield request.getfixturevalue('postgres_session')
		case DatabaseBackend.SQLITE:
			yield request.getfixturevalue('sqlite_session')
		case _:
			raise RuntimeError(f'Unsupported test database backend: {backend!r}')
