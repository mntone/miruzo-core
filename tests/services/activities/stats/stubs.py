from app.models.records import StatsRecord


class StubStatsRepository:
	def __init__(self) -> None:
		self.stats_response: StatsRecord | None = None
		self.upsert_called_with: int | None = None

	def upsert_with_increment(self, ingest_id: int) -> StatsRecord:
		self.upsert_called_with = ingest_id
		if self.stats_response is None:
			raise RuntimeError('stats_response not configured')
		return self.stats_response
