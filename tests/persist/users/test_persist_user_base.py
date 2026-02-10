import sqlite3
from collections.abc import Generator
from typing import Any

import pytest
from sqlalchemy.exc import IntegrityError
from sqlmodel import Session, SQLModel, create_engine

from app.errors import SingletonUserMissingError
from app.models.records import UserRecord
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


def test_create_singleton_if_missing_recovers_from_unique_violation(
	session: Session,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	repo = SQLiteUserRepository(session)
	existing = repo.create_singleton_if_missing()
	existing.daily_love_used = 4
	session.commit()

	original_get = session.get
	get_called_once = False

	def fake_get(model: Any, ident: Any, *args: Any, **kwargs: Any) -> Any:
		nonlocal get_called_once
		if model is UserRecord and ident == 1 and not get_called_once:
			get_called_once = True
			return None
		return original_get(model, ident, *args, **kwargs)

	original_flush = session.flush
	flush_called_once = False

	def fake_flush(*args: Any, **kwargs: Any) -> None:
		nonlocal flush_called_once
		if not flush_called_once:
			flush_called_once = True
			raise IntegrityError(
				'INSERT INTO users (id) VALUES (1)',
				{},
				sqlite3.IntegrityError('UNIQUE constraint failed: users.id'),
			)
		original_flush(*args, **kwargs)

	monkeypatch.setattr(session, 'get', fake_get)
	monkeypatch.setattr(session, 'flush', fake_flush)

	user = repo.create_singleton_if_missing()

	assert user.id == 1
	assert user.daily_love_used == 4


def test_get_singleton_raises_when_missing(session: Session) -> None:
	repo = SQLiteUserRepository(session)

	with pytest.raises(SingletonUserMissingError):
		repo.get_singleton()


def test_increment_daily_love_used_respects_limit(session: Session) -> None:
	repo = SQLiteUserRepository(session)
	user = repo.create_singleton_if_missing()
	user.daily_love_used = 1

	assert repo.increment_daily_love_used(limit=2) is True
	assert repo.increment_daily_love_used(limit=2) is False

	user = repo.get_singleton()
	assert user.daily_love_used == 2


def test_increment_daily_love_used_works_without_singleton_in_identity_map(session: Session) -> None:
	repo = SQLiteUserRepository(session)
	repo.create_singleton_if_missing()
	session.expunge_all()

	assert repo.increment_daily_love_used(limit=2) is True

	user = repo.get_singleton()
	assert user.daily_love_used == 1


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


def test_decrement_daily_love_used_works_without_singleton_in_identity_map(session: Session) -> None:
	repo = SQLiteUserRepository(session)
	user = repo.create_singleton_if_missing()
	user.daily_love_used = 1
	session.flush()
	session.expunge_all()

	assert repo.decrement_daily_love_used() is True

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
