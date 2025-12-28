# pyright: reportAttributeAccessIssue=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false
# pyright: reportUnknownVariableType=false

from abc import ABC, abstractmethod
from typing import TypeVar

from sqlalchemy import Insert, func
from sqlmodel import Session, SQLModel

from app.models.records import StatsRecord

TModel = TypeVar('TModel', bound=SQLModel)


class BaseStatsRepository(ABC):
	def __init__(self, session: Session) -> None:
		self._session = session

	@abstractmethod
	def _build_insert(self, model: type[TModel]) -> Insert: ...

	def upsert_with_increment(
		self,
		ingest_id: int,
	) -> StatsRecord:
		statement = (
			self._build_insert(StatsRecord)
			# INSERT INTO imagestats (ingest_id, view_count) VALUES (:ingest_id, 1)
			.values(
				ingest_id=ingest_id,
				view_count=1,
				last_viewed_at=func.current_timestamp(),
			)
			# ON CONFLICT(ingest_id)
			# DO UPDATE SET view_count = view_count + 1
			.on_conflict_do_update(
				index_elements=['ingest_id'],
				set_={
					'view_count': StatsRecord.view_count + 1,
					'last_viewed_at': func.current_timestamp(),
				},
			)
			# RETURNING ingest_id, hall_of_fame_at, score, view_count, last_viewed_at
			.returning(
				StatsRecord.ingest_id,
				StatsRecord.hall_of_fame_at,
				StatsRecord.score,
				StatsRecord.view_count,
				StatsRecord.last_viewed_at,
			)
		)

		row = self._session.exec(statement).first()
		self._session.commit()

		return StatsRecord(**row._asdict())
