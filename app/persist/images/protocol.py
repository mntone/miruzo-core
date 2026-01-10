from typing import Protocol

from app.models.records import ImageRecord


class ImageRepository(Protocol):
	def select_by_ingest_id(self, ingest_id: int) -> ImageRecord | None: ...
	def insert(self, image: ImageRecord) -> None: ...
