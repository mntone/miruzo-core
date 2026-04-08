from collections.abc import Iterator
from typing import Any

import pytest
from sqlalchemy.orm import Session

from tests.persist.fixtures.backends.runtime import _ensure_runtime_api_available

from app.databases.database import _create_mysql_engine
from app.databases.metadata import metadata

MYSQL_IMAGE = 'mysql:9-oracle'
MYSQL_DB = 'miruzo'
MYSQL_USER = 'm'
MYSQL_PASSWORD = 'miruzo1234'


@pytest.fixture(scope='session')
def mysql_container() -> Iterator[Any]:
	mysql_mod = pytest.importorskip(
		'testcontainers.mysql',
		reason='testcontainers[mysql] is not installed',
	)
	MySqlContainer = mysql_mod.MySqlContainer

	_ensure_runtime_api_available()

	with MySqlContainer(
		image=MYSQL_IMAGE,
		dialect='mysqldb',
		username=MYSQL_USER,
		root_password=MYSQL_PASSWORD,
		password=MYSQL_PASSWORD,
		dbname=MYSQL_DB,
	) as mysql_container:
		yield mysql_container


@pytest.fixture(scope='session')
def mysql_dsn(request: pytest.FixtureRequest) -> Iterator[str]:
	container = request.getfixturevalue('mysql_container')
	yield container.get_connection_url().replace('localhost', '127.0.0.1')


@pytest.fixture()
def mysql_session(request: pytest.FixtureRequest) -> Iterator[Session]:
	dsn = request.getfixturevalue('mysql_dsn')
	engine = _create_mysql_engine(dsn, pool_size=1, max_overflow=2)
	metadata.create_all(engine)
	with Session(engine) as session:
		yield session
	metadata.drop_all(engine)
