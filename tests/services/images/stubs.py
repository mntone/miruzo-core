from datetime import datetime

from app.models.records import ImageRecord


class StubImageRepository:
	def __init__(self) -> None:
		self.list_response: tuple[list[ImageRecord], datetime | None] = ([], None)
		self.list_called_with: dict[str, object] | None = None

		self.one_response: ImageRecord | None = None
		self.one_called_with: int | None = None

	def select_latest(
		self,
		*,
		cursor: datetime | None,
		limit: int,
	) -> tuple[list[ImageRecord], datetime | None]:
		self.list_called_with = {'cursor': cursor, 'limit': limit}
		return self.list_response

	def select_by_ingest_id(self, ingest_id: int) -> ImageRecord | None:
		self.one_called_with = ingest_id
		return self.one_response
