from collections.abc import Iterable
from datetime import datetime
from typing import Callable, final

from tests.stubs.session import StubSession

from app.models.records import StatsRecord
from app.persist.stats.protocol import StatsRepository


@final
class StubStatsRepository:
	def __init__(self) -> None:
		self.stats_list_response: list[StatsRecord] = []

		self.get_one_called_with: int | None = None
		self.create_called_with: int | None = None
		self.create_initial_score: int | None = None
		self.create_response: StatsRecord | None = None

		self.try_set_last_loved_at_called_with: tuple[int, datetime, datetime] | None = None
		self.try_set_last_loved_at_response: bool = True

		self.try_unset_last_loved_at_called_with: tuple[int, datetime] | None = None
		self.try_unset_last_loved_at_response: bool = True

	def get_one(self, ingest_id: int) -> StatsRecord:
		self.get_one_called_with = ingest_id
		for record in self.stats_list_response:
			if record.ingest_id == ingest_id:
				return record
		raise RuntimeError('stats_list_response not configured')

	def create(self, ingest_id: int, *, initial_score: int) -> StatsRecord:
		self.create_called_with = ingest_id
		self.create_initial_score = initial_score
		stats = StatsRecord(
			ingest_id=ingest_id,
			score=initial_score,
			score_evaluated=initial_score,
		)
		self.create_response = stats
		self.stats_list_response.append(stats)
		return stats

	def try_set_last_loved_at(
		self,
		ingest_id: int,
		*,
		last_loved_at: datetime,
		since_occurred_at: datetime,
	) -> bool:
		self.try_set_last_loved_at_called_with = (ingest_id, last_loved_at, since_occurred_at)
		return self.try_set_last_loved_at_response

	def try_unset_last_loved_at(
		self,
		ingest_id: int,
		*,
		since_occurred_at: datetime,
	) -> bool:
		self.try_unset_last_loved_at_called_with = (ingest_id, since_occurred_at)
		return self.try_unset_last_loved_at_response

	def iterable(self) -> Iterable[StatsRecord]:
		return sorted(self.stats_list_response, key=lambda record: record.ingest_id)


def create_stub_stats_repository_factory(repo: StatsRepository) -> Callable[[StubSession], StatsRepository]:
	def factory(_: StubSession) -> StatsRepository:
		return repo

	return factory
