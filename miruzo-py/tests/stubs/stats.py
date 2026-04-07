from typing import final

from app.persist.stats.protocol import StatsCreateInput


@final
class StubStatsRepository:
	def __init__(self) -> None:
		self.create_called_with: StatsCreateInput | None = None

	def create(self, entry: StatsCreateInput) -> None:
		self.create_called_with = entry
