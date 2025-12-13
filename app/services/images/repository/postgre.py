from sqlalchemy.dialects.postgresql import insert
from sqlmodel import select

from app.models.records import ImageRecord, StatsRecord
from app.services.images.repository.base import ImageRepository


class PostgreSQLImageRepository(ImageRepository):
	def get_detail_with_stats(
		self,
		image_id: int,
	) -> tuple[ImageRecord, StatsRecord] | None:
		row = self._session.exec(
			select(ImageRecord, StatsRecord)
			.join(StatsRecord, StatsRecord.image_id == ImageRecord.id)
			.where(ImageRecord.id == image_id),
		).first()

		if row is None:
			return None

		return row

	def upsert_stats_with_increment(self, image_id: int) -> StatsRecord:
		return self._upsert_stats_with_increment(insert, image_id)
