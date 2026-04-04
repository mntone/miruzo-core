from typing import Protocol

from app.models.records import ImageRecord


class ImageRepository(Protocol):
	def insert(self, image: ImageRecord) -> None: ...
