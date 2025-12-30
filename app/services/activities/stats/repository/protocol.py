from typing import Protocol

from app.models.records import StatsRecord


class StatsRepository(Protocol):
	def get_or_create(
		self,
		ingest_id: int,
		*,
		initial_score: int,
	) -> StatsRecord: ...
