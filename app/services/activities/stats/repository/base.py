# pyright: reportAttributeAccessIssue=false
# pyright: reportArgumentType=false
# pyright: reportOptionalMemberAccess=false
# pyright: reportOptionalOperand=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false
# pyright: reportUnknownVariableType=false

from abc import ABC, abstractmethod
from collections.abc import Iterable
from datetime import datetime
from typing import TypeVar

from sqlalchemy import Insert, or_, true, update
from sqlalchemy.exc import IntegrityError
from sqlmodel import Session, SQLModel, select

from app.models.records import StatsRecord

TModel = TypeVar('TModel', bound=SQLModel)


class BaseStatsRepository(ABC):
	def __init__(self, session: Session) -> None:
		self._session = session

	@abstractmethod
	def _is_unique_violation(self, error: IntegrityError) -> bool: ...

	def get_one(self, ingest_id: int) -> StatsRecord:
		stats = self._session.get_one(StatsRecord, ingest_id)

		return stats

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

	def try_set_last_loved_at(
		self,
		ingest_id: int,
		*,
		last_loved_at: datetime,
		since_occurred_at: datetime,
	) -> bool:
		statement = (
			update(StatsRecord)
			.where(StatsRecord.ingest_id == ingest_id)
			.where(
				or_(
					StatsRecord.last_loved_at.is_(None),
					StatsRecord.last_loved_at < since_occurred_at,
				),
			)
			.values(last_loved_at=last_loved_at)
		)

		result = self._session.exec(statement)

		return result.rowcount == 1

	def try_unset_last_loved_at(
		self,
		ingest_id: int,
		*,
		since_occurred_at: datetime,
	) -> bool:
		statement = (
			update(StatsRecord)
			.where(StatsRecord.ingest_id == ingest_id)
			.where(
				StatsRecord.last_loved_at.is_not(None),
				StatsRecord.last_loved_at >= since_occurred_at,
			)
			.values(last_loved_at=None)
		)

		result = self._session.exec(statement)

		return result.rowcount == 1

	@abstractmethod
	def _build_insert(self, model: type[TModel]) -> Insert: ...

	def iterable(self) -> Iterable[StatsRecord]:
		last_ingest_id = None

		while True:
			statement = (
				select(StatsRecord)
				.where(StatsRecord.ingest_id > last_ingest_id if last_ingest_id is not None else true())
				.order_by(StatsRecord.ingest_id.asc())
				.limit(500)
			)

			rows = self._session.exec(statement).all()
			if not rows:
				break

			yield from rows

			last_ingest_id = rows[-1].ingest_id
