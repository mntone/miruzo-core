from sqlmodel import Session

from app.persist.stats.base import BaseStatsRepository
from app.persist.stats.protocol import StatsRepository


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

	return BaseStatsRepository(session)
