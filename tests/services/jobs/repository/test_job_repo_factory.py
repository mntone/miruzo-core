from typing import Any, Generator

import pytest
from sqlmodel import Session, create_engine

from app.config.environments import DatabaseBackend, env
from app.services.jobs.repository.factory import create_job_repository
from app.services.jobs.repository.postgre import PostgreSQLJobRepository
from app.services.jobs.repository.sqlite import SQLiteJobRepository


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	with Session(engine) as session:
		yield session


def test_create_job_repository_uses_sqlite(session: Session, monkeypatch: pytest.MonkeyPatch) -> None:
	monkeypatch.setattr(env, 'database_backend', DatabaseBackend.SQLITE)

	repo = create_job_repository(session)

	assert isinstance(repo, SQLiteJobRepository)


def test_create_job_repository_uses_postgres(session: Session, monkeypatch: pytest.MonkeyPatch) -> None:
	monkeypatch.setattr(env, 'database_backend', DatabaseBackend.POSTGRE_SQL)

	repo = create_job_repository(session)

	assert isinstance(repo, PostgreSQLJobRepository)


def test_create_job_repository_rejects_unknown_backend(
	session: Session,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	monkeypatch.setattr(env, 'database_backend', 'mysql')

	with pytest.raises(ValueError, match='Unsupported database type'):
		create_job_repository(session)
