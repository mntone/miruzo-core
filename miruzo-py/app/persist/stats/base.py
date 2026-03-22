# pyright: reportAttributeAccessIssue=false
# pyright: reportArgumentType=false
# pyright: reportOptionalMemberAccess=false
# pyright: reportOptionalOperand=false
# pyright: reportUnknownVariableType=false

from collections.abc import Iterable
from typing import TypeVar, final

from sqlalchemy import true
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
