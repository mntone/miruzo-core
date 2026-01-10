# pyright: reportAttributeAccessIssue=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false

from collections.abc import Collection, Sequence
from datetime import datetime
from typing import final

from sqlmodel import Session, select

from app.models.enums import ActionKind
from app.models.records import ActionRecord


@final
class BaseActionRepository:
	def __init__(self, session: Session) -> None:
		self._session = session

	def select_by_ingest_id(self, ingest_id: int) -> Sequence[ActionRecord]:
		statement = select(ActionRecord)

		statement = statement.where(ActionRecord.ingest_id == ingest_id)

		statement = statement.order_by(ActionRecord.occurred_at.asc())

		items = self._session.exec(statement).all()

		return items

	def select_latest_one(
		self,
		ingest_id: int,
		*,
		kind: ActionKind,
		since_occurred_at: datetime,
		until_occurred_at: datetime | None = None,
		require_unique: bool = False,
	) -> ActionRecord | None:
		"""
		Return the latest matching ActionRecord.
		Time range is interpreted as [since_occurred_at, until_occurred_at).
		When require_unique is True, raise if multiple rows exist.
		"""

		statement = (
			select(ActionRecord)
			.where(
				ActionRecord.ingest_id == ingest_id,
				ActionRecord.kind == kind,
				ActionRecord.occurred_at >= since_occurred_at,
			)
			.order_by(
				ActionRecord.occurred_at.desc(),
				ActionRecord.id.desc(),
			)
			.limit(2 if require_unique else 1)
		)

		if until_occurred_at is not None:
			statement = statement.where(
				ActionRecord.occurred_at < until_occurred_at,
			)

		row = self._session.exec(statement).one_or_none()

		return row

	def select_latest_one_by_multiple_kinds(
		self,
		ingest_id: int,
		*,
		kinds: Collection[ActionKind],
		since_occurred_at: datetime | None = None,
		until_occurred_at: datetime | None = None,
		require_unique: bool = False,
	) -> ActionRecord | None:
		"""
		Return the latest matching ActionRecord.
		Time range is interpreted as [since_occurred_at, until_occurred_at).
		When require_unique is True, raise if multiple rows exist.
		"""

		if not kinds:
			return None

		statement = (
			select(ActionRecord)
			.where(
				ActionRecord.ingest_id == ingest_id,
				ActionRecord.kind.in_(kinds),
			)
			.order_by(
				ActionRecord.occurred_at.desc(),
				ActionRecord.id.desc(),
			)
			.limit(2 if require_unique else 1)
		)

		if since_occurred_at is not None:
			statement = statement.where(
				ActionRecord.occurred_at >= since_occurred_at,
			)

		if until_occurred_at is not None:
			statement = statement.where(
				ActionRecord.occurred_at < until_occurred_at,
			)

		row = self._session.exec(statement).one_or_none()

		return row

	def insert(
		self,
		ingest_id: int,
		*,
		kind: ActionKind,
		occurred_at: datetime,
	) -> ActionRecord:
		action = ActionRecord(
			ingest_id=ingest_id,
			kind=kind,
			occurred_at=occurred_at,
		)

		self._session.add(action)
		self._session.flush()
		self._session.refresh(action)

		return action
