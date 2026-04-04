from typing import TypeVar, final

from sqlmodel import Session, SQLModel

from app.models.records import StatsRecord

TModel = TypeVar('TModel', bound=SQLModel)


@final
class BaseStatsRepository:
	def __init__(self, session: Session) -> None:
		self._session = session

	def get_one(self, ingest_id: int) -> StatsRecord:
		stats = self._session.get_one(StatsRecord, ingest_id)

		return stats

	def create(
		self,
		ingest_id: int,
		*,
		initial_score: int,
	) -> StatsRecord:
		stats = StatsRecord(
			ingest_id=ingest_id,
			score=initial_score,
			score_evaluated=initial_score,
		)
		self._session.add(stats)
		self._session.flush()
		self._session.refresh(stats)
		return stats
