from collections.abc import Iterable
from datetime import datetime
from typing import Protocol

from app.models.records import StatsRecord


class StatsRepository(Protocol):
	def get_one(self, ingest_id: int) -> StatsRecord: ...

	def get_or_create(
		self,
		ingest_id: int,
		*,
		initial_score: int,
	) -> StatsRecord: ...

	def try_set_last_loved_at(
		self,
		ingest_id: int,
		*,
		last_loved_at: datetime,
		since_occurred_at: datetime,
	) -> bool: ...

	def try_unset_last_loved_at(
		self,
		ingest_id: int,
		*,
		since_occurred_at: datetime,
	) -> bool: ...

	def iterable(self) -> Iterable[StatsRecord]:
		"""Yield StatsRecord entries in ingest_id order using batch paging."""
		...
