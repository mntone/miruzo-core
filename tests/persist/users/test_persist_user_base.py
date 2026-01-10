from collections.abc import Generator
from typing import Any

import pytest
from sqlmodel import Session, SQLModel, create_engine

from app.persist.users.sqlite import SQLiteUserRepository


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	SQLModel.metadata.create_all(engine)
	with Session(engine) as session:
		yield session


def test_get_or_create_singleton_returns_existing(session: Session) -> None:
	repo = SQLiteUserRepository(session)

	user = repo.get_or_create_singleton()

	user_again = repo.get_or_create_singleton()

	assert user_again is user
