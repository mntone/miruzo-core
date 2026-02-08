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


def test_try_increment_daily_love_used_respects_limit(session: Session) -> None:
	repo = SQLiteUserRepository(session)

	assert repo.try_increment_daily_love_used(limit=2) is True
	assert repo.try_increment_daily_love_used(limit=2) is True
	assert repo.try_increment_daily_love_used(limit=2) is False

	user = repo.get_or_create_singleton()
	assert user.daily_love_used == 2


def test_try_decrement_daily_love_used_stops_at_zero(session: Session) -> None:
	repo = SQLiteUserRepository(session)
	user = repo.get_or_create_singleton()
	user.daily_love_used = 1

	assert repo.try_decrement_daily_love_used() is True
	assert repo.try_decrement_daily_love_used() is False

	user = repo.get_or_create_singleton()
	assert user.daily_love_used == 0


def test_reset_daily_love_used_sets_zero(session: Session) -> None:
	repo = SQLiteUserRepository(session)
	user = repo.get_or_create_singleton()
	user.daily_love_used = 5

	repo.reset_daily_love_used()

	user = repo.get_or_create_singleton()
	assert user.daily_love_used == 0
