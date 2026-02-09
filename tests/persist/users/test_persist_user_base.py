from collections.abc import Generator
from typing import Any

import pytest
from sqlmodel import Session, SQLModel, create_engine

from app.errors import SingletonUserMissingError
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


def test_create_singleton_if_missing_returns_existing(session: Session) -> None:
	repo = SQLiteUserRepository(session)

	user = repo.create_singleton_if_missing()

	user_again = repo.create_singleton_if_missing()

	assert user_again is user


def test_get_singleton_raises_when_missing(session: Session) -> None:
	repo = SQLiteUserRepository(session)

	with pytest.raises(SingletonUserMissingError):
		repo.get_singleton()


def test_increment_daily_love_used_respects_limit(session: Session) -> None:
	repo = SQLiteUserRepository(session)
	repo.create_singleton_if_missing()

	assert repo.increment_daily_love_used(limit=2) is True
	assert repo.increment_daily_love_used(limit=2) is True
	assert repo.increment_daily_love_used(limit=2) is False

	user = repo.get_singleton()
	assert user.daily_love_used == 2


def test_increment_daily_love_used_raises_when_missing(session: Session) -> None:
	repo = SQLiteUserRepository(session)

	with pytest.raises(SingletonUserMissingError):
		repo.increment_daily_love_used(limit=1)


def test_decrement_daily_love_used_stops_at_zero(session: Session) -> None:
	repo = SQLiteUserRepository(session)
	user = repo.create_singleton_if_missing()
	user.daily_love_used = 1

	assert repo.decrement_daily_love_used() is True
	assert repo.decrement_daily_love_used() is False

	user = repo.get_singleton()
	assert user.daily_love_used == 0


def test_decrement_daily_love_used_raises_when_missing(session: Session) -> None:
	repo = SQLiteUserRepository(session)

	with pytest.raises(SingletonUserMissingError):
		repo.decrement_daily_love_used()


def test_reset_daily_love_used_sets_zero(session: Session) -> None:
	repo = SQLiteUserRepository(session)
	user = repo.create_singleton_if_missing()
	user.daily_love_used = 5

	repo.reset_daily_love_used()

	user = repo.get_singleton()
	assert user.daily_love_used == 0


def test_reset_daily_love_used_raises_when_missing(session: Session) -> None:
	repo = SQLiteUserRepository(session)

	with pytest.raises(SingletonUserMissingError):
		repo.reset_daily_love_used()
