from typing import Protocol

from app.models.records import StatsRecord


class StatsRepository(Protocol):
	def upsert_with_increment(
		self,
		ingest_id: int,
	) -> StatsRecord:
		"""Increment views (upserting as needed) and return the latest row."""
		...
