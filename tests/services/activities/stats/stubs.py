from collections.abc import Iterable, Sequence
from typing import final

from app.models.records import StatsRecord


@final
class StubStatsRepository:
	def __init__(self) -> None:
		self.stats_response: StatsRecord | None = None
		self.get_or_create_called_with: int | None = None
		self.get_or_create_initial_score: int | None = None

		self.stats_list_response: Sequence[StatsRecord] = []

	def get_or_create(self, ingest_id: int, *, initial_score: int) -> StatsRecord:
		self.get_or_create_called_with = ingest_id
		self.get_or_create_initial_score = initial_score
		if self.stats_response is None:
			raise RuntimeError('stats_response not configured')
		return self.stats_response

	def iterable(self) -> Iterable[StatsRecord]:
		return self.stats_list_response
