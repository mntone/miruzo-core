from collections.abc import Iterator

import pytest
from sqlalchemy.orm import Session

from app.databases.database import _create_sqlite3_engine
from app.databases.metadata import metadata


@pytest.fixture()
def sqlite_session() -> Iterator[Session]:
	engine = _create_sqlite3_engine('sqlite+pysqlite:///:memory:')
	metadata.create_all(engine)
	with Session(engine) as session:
		yield session
