from collections.abc import Collection, Sequence
from datetime import datetime
from typing import Protocol

from app.models.enums import ActionKind
from app.models.records import ActionRecord


class ActionRepository(Protocol):
	def select_by_ingest_id(self, ingest_id: int) -> Sequence[ActionRecord]: ...

	def select_latest_one(
		self,
		ingest_id: int,
		*,
		kind: ActionKind,
		since_occurred_at: datetime,
		until_occurred_at: datetime | None = None,
		require_unique: bool = False,
	) -> ActionRecord | None: ...

	def select_latest_one_by_multiple_kinds(
		self,
		ingest_id: int,
		*,
		kinds: Collection[ActionKind],
		since_occurred_at: datetime | None = None,
		until_occurred_at: datetime | None = None,
		require_unique: bool = False,
	) -> ActionRecord | None: ...

	def select_latest_effective_love(
		self,
		ingest_id: int,
		*,
		since_occurred_at: datetime | None = None,
		until_occurred_at: datetime | None = None,
	) -> ActionRecord | None: ...

	def insert(
		self,
		ingest_id: int,
		*,
		kind: ActionKind,
		occurred_at: datetime,
	) -> ActionRecord: ...
