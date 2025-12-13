from sqlalchemy.dialects.sqlite import insert

from app.models.records import ImageRecord, StatsRecord
from app.services.images.repository.base import ImageRepository


class SQLiteImageRepository(ImageRepository):
	def get_detail_with_stats(
		self,
		image_id: int,
	) -> tuple[ImageRecord, StatsRecord] | None:
		"""
		Fetch an image and its stats in a single request.

		Returns:
			tuple of image record and stats if both exist, otherwise ``None``.
		"""
		image = self._session.get(ImageRecord, image_id)

		if image is None:
			return None

		stats = self._session.get(StatsRecord, image_id)

		return image, stats

	def upsert_stats_with_increment(self, image_id: int) -> StatsRecord:
		return self._upsert_stats_with_increment(insert, image_id)
