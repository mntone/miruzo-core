from typing import final

from app.models.records import StatsRecord
from app.services.activities.stats.repository.protocol import StatsRepository


@final
class StatsService:
	def __init__(self, repository: StatsRepository) -> None:
		self._repository = repository

	def get_by_ingest_id(self, ingest_id: int) -> StatsRecord:
		record = self._repository.upsert_with_increment(ingest_id)

		return record
