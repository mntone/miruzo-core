from collections.abc import Iterator
from typing import cast

import pytest
from sqlalchemy.orm import Session
from testcontainers.postgres import PostgresContainer

from tests.persist.fixtures.backends.runtime import _ensure_runtime_api_available

from app.databases.database import _create_postgres_engine
from app.databases.metadata import metadata

POSTGRES_IMAGE = 'postgres:18-alpine'
POSTGRES_DB = 'miruzo'
POSTGRES_USER = 'm'
POSTGRES_PASSWORD = 'miruzo1234'
POSTGRES_ARGS = '--encoding=UTF8 --lc-collate=C --lc-ctype=C'


@pytest.fixture(scope='session')
def postgres_container() -> Iterator[PostgresContainer]:
	_ensure_runtime_api_available()

	with PostgresContainer(
		image=POSTGRES_IMAGE,
		username=POSTGRES_USER,
		password=POSTGRES_PASSWORD,
		dbname=POSTGRES_DB,
		driver='psycopg2',
	).with_env(
		'POSTGRES_INITDB_ARGS',
		POSTGRES_ARGS,
	) as postgres_container:
		yield postgres_container


@pytest.fixture(scope='session')
def postgres_dsn(request: pytest.FixtureRequest) -> Iterator[str]:
	container = cast(PostgresContainer, request.getfixturevalue('postgres_container'))
	yield container.get_connection_url()


@pytest.fixture()
def postgres_session(request: pytest.FixtureRequest) -> Iterator[Session]:
	dsn = request.getfixturevalue('postgres_dsn')
	engine = _create_postgres_engine(dsn, pool_size=1, max_overflow=2)
	metadata.create_all(engine)
	with Session(engine) as postgres_session:
		yield postgres_session
	metadata.drop_all(engine)
