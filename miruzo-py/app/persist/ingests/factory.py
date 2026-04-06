from sqlalchemy.orm import Session

from app.config.environments import DatabaseBackend, env
from app.persist.ingests.base import _IngestRepositoryBaseImpl
from app.persist.ingests.postgres import _IngestRepositoryPostgresImpl
from app.persist.ingests.protocol import IngestRepository


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

	if env.database_backend == DatabaseBackend.SQLITE:
		return _IngestRepositoryBaseImpl(session)
	elif env.database_backend == DatabaseBackend.POSTGRE_SQL:
		return _IngestRepositoryPostgresImpl(session)
	else:
		raise ValueError(f'Unsupported database type: {env.database_backend}')
