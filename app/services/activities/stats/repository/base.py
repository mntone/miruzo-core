# pyright: reportAttributeAccessIssue=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false
# pyright: reportUnknownVariableType=false

from abc import ABC, abstractmethod
from typing import TypeVar

from sqlalchemy import Insert
from sqlalchemy.exc import IntegrityError
from sqlmodel import Session, SQLModel

from app.models.records import StatsRecord

TModel = TypeVar('TModel', bound=SQLModel)


class BaseStatsRepository(ABC):
	def __init__(self, session: Session) -> None:
		self._session = session

	@abstractmethod
	def _is_unique_violation(self, error: IntegrityError) -> bool: ...

	def get_or_create(
		self,
		ingest_id: int,
		*,
		initial_score: int,
	) -> StatsRecord:
		stats = self._session.get(StatsRecord, ingest_id)
		if stats is not None:
			return stats

		stats = StatsRecord(
			ingest_id=ingest_id,
			score=initial_score,
			view_count=0,
			last_viewed_at=None,
			hall_of_fame_at=None,
		)
		self._session.add(stats)

		try:
			self._session.flush()
		except IntegrityError as exc:
			self._session.rollback()
			if not self._is_unique_violation(exc):
				raise
			stats = self._session.get_one(StatsRecord, ingest_id)

		return stats

	@abstractmethod
	def _build_insert(self, model: type[TModel]) -> Insert: ...
