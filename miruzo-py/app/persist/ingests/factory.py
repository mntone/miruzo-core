from sqlalchemy.orm import Session

from app.config.environments import DatabaseBackend, env
from app.models.ingest import MAX_EXECUTIONS
from app.persist.ingests.base import _IngestRepositoryBaseImpl
from app.persist.ingests.postgres import _IngestRepositoryPostgresImpl
from app.persist.ingests.protocol import IngestRepository


def _create_ingest_repository_from_backend(
	session: Session,
	*,
	backend: DatabaseBackend,
	max_executions: int,
) -> IngestRepository:
	match backend:
		case DatabaseBackend.POSTGRE_SQL:
			return _IngestRepositoryPostgresImpl(session, max_executions=max_executions)
		case DatabaseBackend.SQLITE:
			return _IngestRepositoryBaseImpl(session, max_executions=max_executions)
		case _:
			raise ValueError(f'Unsupported database type: {backend}')


def create_ingest_repository(session: Session) -> IngestRepository:
	"""
	Build a ingest repository implementation for the configured backend.

	Args:
		session: SQLAlchemy session bound to the current database engine.

	Returns:
		Concrete repository tied to the active backend.

	Raises:
		ValueError: if the configured backend is unsupported.
	"""

	return _create_ingest_repository_from_backend(
		session,
		backend=env.database_backend,
		max_executions=MAX_EXECUTIONS,
	)
