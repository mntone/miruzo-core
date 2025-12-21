# pyright: reportArgumentType=false
# pyright: reportUnknownMemberType=false

from sqlalchemy.dialects.postgresql import Insert as PostgreInsert
from sqlalchemy.dialects.postgresql import insert as postgre_insert
from sqlmodel import SQLModel, select

from app.models.records import ImageRecord, StatsRecord
from app.services.images.repository.base import ImageRepository


class PostgreSQLImageRepository(ImageRepository):
	def get_detail_with_stats(
		self,
		ingest_id: int,
	) -> tuple[ImageRecord, StatsRecord | None] | None:
		row = self._session.exec(
			select(ImageRecord, StatsRecord)
			.join(StatsRecord, StatsRecord.ingest_id == ImageRecord.ingest_id)
			.where(ImageRecord.ingest_id == ingest_id),
		).first()

		if row is None:
			return None

		return row

	def _build_insert(self, model: type[SQLModel]) -> PostgreInsert:
		return postgre_insert(model.__table__)  # pyright: ignore[reportAttributeAccessIssue, reportUnknownArgumentType]
