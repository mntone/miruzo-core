# pyright: reportAttributeAccessIssue=false
# pyright: reportArgumentType=false
# pyright: reportOptionalMemberAccess=false
# pyright: reportOptionalOperand=false
# pyright: reportUnknownVariableType=false

from collections.abc import Iterable
from datetime import datetime
from typing import TypeVar, final

from sqlalchemy import or_, true, update
from sqlmodel import Session, SQLModel, select

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
