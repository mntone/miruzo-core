from typing import Any, Generator

import pytest
from sqlalchemy import create_engine
from sqlalchemy.orm import Session

from app.databases.metadata import metadata


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	metadata.create_all(engine)
	with Session(engine) as session:
		yield session
