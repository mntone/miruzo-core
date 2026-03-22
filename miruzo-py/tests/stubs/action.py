from collections.abc import Collection, Sequence
from dataclasses import dataclass
from datetime import datetime
from typing import final

from sqlalchemy.exc import MultipleResultsFound

from app.models.enums import ActionKind
from app.models.records import ActionRecord


@dataclass(frozen=True, slots=True)
@final
class _StubActionRepository_SelectOneArgs:
	ingest_id: int
	kinds: Collection[ActionKind]
	since_occurred_at: datetime | None
	until_occurred_at: datetime | None
	require_unique: bool


@dataclass(frozen=True, slots=True)
@final
class _StubActionRepository_InsertArgs:
	ingest_id: int
	kind: ActionKind
	occurred_at: datetime


@final
class StubActionRepository:
	def __init__(self) -> None:
		self.select_called_with: int | None = None
		self.select_one_called_with: _StubActionRepository_SelectOneArgs | None = None
		self.insert_called_with: list[_StubActionRepository_InsertArgs] = []
		self.actions: list[ActionRecord] = []

	def select_by_ingest_id(self, ingest_id: int) -> Sequence[ActionRecord]:
		self.select_called_with = ingest_id
		items = [action for action in self.actions if action.ingest_id == ingest_id]
		return sorted(items, key=self._sort_key)

	def select_latest_one(
		self,
		ingest_id: int,
		*,
		kind: ActionKind,
		since_occurred_at: datetime,
		until_occurred_at: datetime | None = None,
		require_unique: bool = False,
	) -> ActionRecord | None:
		return self.select_latest_one_by_multiple_kinds(
			ingest_id,
			kinds=(kind,),
			since_occurred_at=since_occurred_at,
			until_occurred_at=until_occurred_at,
			require_unique=require_unique,
		)

	def select_latest_one_by_multiple_kinds(
		self,
		ingest_id: int,
		*,
		kinds: Collection[ActionKind],
		since_occurred_at: datetime | None = None,
		until_occurred_at: datetime | None = None,
		require_unique: bool = False,
	) -> ActionRecord | None:
		self.select_one_called_with = _StubActionRepository_SelectOneArgs(
			ingest_id=ingest_id,
			kinds=kinds,
			since_occurred_at=since_occurred_at,
			until_occurred_at=until_occurred_at,
			require_unique=require_unique,
		)

		candidates: list[ActionRecord] = []
		for action in self.actions:
			if action.ingest_id != ingest_id:
				continue

			if action.kind not in kinds:
				continue

			if since_occurred_at and action.occurred_at < since_occurred_at:
				continue

			if until_occurred_at and action.occurred_at >= until_occurred_at:
				continue

			candidates.append(action)

		if not candidates:
			return None

		candidates = sorted(candidates, key=self._sort_key, reverse=True)
		if require_unique and len(candidates) > 1:
			raise MultipleResultsFound('multiple action')

		return candidates[0]

	def insert(
		self,
		ingest_id: int,
		*,
		kind: ActionKind,
		occurred_at: datetime,
	) -> ActionRecord:
		self.insert_called_with.append(
			_StubActionRepository_InsertArgs(
				ingest_id=ingest_id,
				kind=kind,
				occurred_at=occurred_at,
			),
		)

		action = ActionRecord(
			id=1,
			ingest_id=ingest_id,
			kind=kind,
			occurred_at=occurred_at,
		)

		self.actions.append(action)
		return action

	@staticmethod
	def _sort_key(action: ActionRecord) -> tuple[datetime, int]:
		return (action.occurred_at, action.id or 0)
