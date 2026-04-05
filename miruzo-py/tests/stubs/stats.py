from typing import Callable, final

from tests.stubs.session import StubSession

from app.persist.stats.protocol import StatsCreateInput, StatsRepository


@final
class StubStatsRepository:
	def __init__(self) -> None:
		self.create_called_with: StatsCreateInput | None = None

	def create(self, entry: StatsCreateInput) -> None:
		self.create_called_with = entry


def create_stub_stats_repository_factory(repo: StatsRepository) -> Callable[[StubSession], StatsRepository]:
	def factory(_: StubSession) -> StatsRepository:
		return repo

	return factory
