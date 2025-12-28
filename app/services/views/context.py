from typing import final

from app.models.api.images.responses import ContextResponse
from app.services.activities.stats.service import StatsService
from app.services.images.query import ImageQueryService


@final
class ContextService:
	def __init__(self, *, image_query: ImageQueryService, stats: StatsService) -> None:
		self._image = image_query
		self._stats = stats

	def get_context(
		self,
		ingest_id: int,
	) -> ContextResponse | None:
		"""
		Return a single image detail payload.

		Fetches the image record, increments view stats, and normalizes
		variant layers based on allowed formats.
		"""
		image = self._image.get_by_ingest_id(ingest_id)

		if image is None:
			return None

		stats = self._stats.get_by_ingest_id(ingest_id)

		response = ContextResponse.from_record(image, stats)

		return response
