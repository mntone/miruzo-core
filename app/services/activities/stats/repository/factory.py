from sqlmodel import Session

from app.config.environments import DatabaseBackend, env
from app.services.activities.stats.repository.postgre import PostgreSQLStatsRepository
from app.services.activities.stats.repository.protocol import StatsRepository
from app.services.activities.stats.repository.sqlite import SQLiteStatsRepository


def create_stats_repository(session: Session) -> StatsRepository:
	"""
	Build an image repository implementation for the configured backend.

	Args:
		session: SQLModel session bound to the current database engine.

	Returns:
		Concrete repository tied to the active backend.

	Raises:
		ValueError: if the configured backend is unsupported.
	"""
	if env.database_backend == DatabaseBackend.SQLITE:
		return SQLiteStatsRepository(session)
	elif env.database_backend == DatabaseBackend.POSTGRE_SQL:
		return PostgreSQLStatsRepository(session)
	else:
		raise ValueError(f'Unsupported database type: {env.database_backend}')
