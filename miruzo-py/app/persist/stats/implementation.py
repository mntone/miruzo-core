from typing import final

from sqlalchemy import insert
from sqlalchemy.orm import Session

from app.databases.tables import stats_table
from app.persist.stats.protocol import StatsCreateInput, StatsRepository


@final
class _StatsRepositoryImpl:
	def __init__(self, session: Session) -> None:
		self._session = session

	def create(self, entry: StatsCreateInput) -> None:
		stmt = insert(stats_table).values(
			ingest_id=entry.ingest_id,
			score=entry.initial_score,
			score_evaluated=entry.initial_score,
		)
		self._session.execute(stmt)


def create_stats_repository(session: Session) -> StatsRepository:
	"""
	Build an stats repository implementation for the configured backend.

	Args:
		session: SQLAlchemy session bound to the current database engine.

	Returns:
		Concrete repository tied to the active backend.
	"""

	return _StatsRepositoryImpl(session)
