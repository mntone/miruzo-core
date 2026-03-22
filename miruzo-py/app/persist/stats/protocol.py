from collections.abc import Iterable
from typing import Protocol

from app.models.records import StatsRecord


class StatsRepository(Protocol):
	def get_one(self, ingest_id: int) -> StatsRecord: ...

	def create(
		self,
		ingest_id: int,
		*,
		initial_score: int,
	) -> StatsRecord: ...

	def iterable(self) -> Iterable[StatsRecord]:
		"""Yield StatsRecord entries in ingest_id order using batch paging."""
		...
