from typing import Any, Generator

import pytest
from sqlmodel import Session, create_engine

from app.config.environments import DatabaseBackend, env
from app.persist.users.factory import create_user_repository
from app.persist.users.postgre import PostgreSQLUserRepository
from app.persist.users.sqlite import SQLiteUserRepository


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	with Session(engine) as session:
		yield session


def test_create_user_repository_uses_sqlite(session: Session, monkeypatch: pytest.MonkeyPatch) -> None:
	monkeypatch.setattr(env, 'database_backend', DatabaseBackend.SQLITE)

	repo = create_user_repository(session)

	assert isinstance(repo, SQLiteUserRepository)


def test_create_user_repository_uses_postgres(session: Session, monkeypatch: pytest.MonkeyPatch) -> None:
	monkeypatch.setattr(env, 'database_backend', DatabaseBackend.POSTGRE_SQL)

	repo = create_user_repository(session)

	assert isinstance(repo, PostgreSQLUserRepository)


def test_create_user_repository_rejects_unknown_backend(
	session: Session,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	monkeypatch.setattr(env, 'database_backend', 'mysql')

	with pytest.raises(ValueError, match='Unsupported database type'):
		create_user_repository(session)
