import pytest
from sqlalchemy.orm import Session

from app.config.environments import DatabaseBackend, env
from app.persist.ingests.base import _IngestRepositoryBaseImpl
from app.persist.ingests.factory import create_ingest_repository
from app.persist.ingests.postgres import _IngestRepositoryPostgresImpl


def test_create_ingest_repository_uses_sqlite(session: Session, monkeypatch: pytest.MonkeyPatch) -> None:
	monkeypatch.setattr(env, 'database_backend', DatabaseBackend.SQLITE)

	repo = create_ingest_repository(session)
	assert isinstance(repo, _IngestRepositoryBaseImpl)


def test_create_ingest_repository_uses_postgres(session: Session, monkeypatch: pytest.MonkeyPatch) -> None:
	monkeypatch.setattr(env, 'database_backend', DatabaseBackend.POSTGRE_SQL)

	repo = create_ingest_repository(session)
	assert isinstance(repo, _IngestRepositoryPostgresImpl)


def test_create_ingest_repository_rejects_unknown_backend(
	session: Session,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	monkeypatch.setattr(env, 'database_backend', 'mysql')

	with pytest.raises(ValueError, match='Unsupported database type'):
		create_ingest_repository(session)
