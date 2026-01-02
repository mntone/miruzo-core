from collections.abc import Sequence
from dataclasses import dataclass
from datetime import datetime
from typing import final

from app.models.enums import ActionKind
from app.models.records import ActionRecord


@dataclass(frozen=True, slots=True)
@final
class _StubActionRepository_SelectOneArgs:
	ingest_id: int
	kind: ActionKind
	since_occurred_at: datetime
	until_occurred_at: datetime


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
		return [action for action in self.actions if action.ingest_id == ingest_id]

	def select_one_by(
		self,
		ingest_id: int,
		*,
		kind: ActionKind,
		since_occurred_at: datetime,
		until_occurred_at: datetime,
	) -> ActionRecord | None:
		self.select_one_called_with = _StubActionRepository_SelectOneArgs(
			ingest_id=ingest_id,
			kind=kind,
			since_occurred_at=since_occurred_at,
			until_occurred_at=until_occurred_at,
		)

		target_action: ActionRecord | None = None
		for action in self.actions:
			if action.ingest_id != ingest_id:
				continue

			if action.kind != kind:
				continue

			if since_occurred_at > action.occurred_at:
				continue

			if action.occurred_at >= until_occurred_at:
				continue

			target_action = action
			break

		return target_action

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
