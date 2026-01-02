from sqlmodel import Session

from app.config.environments import DatabaseBackend, env
from app.services.jobs.repository.postgre import PostgreSQLJobRepository
from app.services.jobs.repository.protocol import JobRepository
from app.services.jobs.repository.sqlite import SQLiteJobRepository


def create_job_repository(session: Session) -> JobRepository:
	"""
	Build a job repository implementation for the configured backend.

	Args:
		session: SQLModel session bound to the current database engine.

	Returns:
		Concrete repository tied to the active backend.

	Raises:
		ValueError: if the configured backend is unsupported.
	"""

	if env.database_backend == DatabaseBackend.SQLITE:
		return SQLiteJobRepository(session)
	elif env.database_backend == DatabaseBackend.POSTGRE_SQL:
		return PostgreSQLJobRepository(session)
	else:
		raise ValueError(f'Unsupported database type: {env.database_backend}')
