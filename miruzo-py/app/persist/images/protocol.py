from typing import Protocol

from app.models.image import Image


class ImageRepository(Protocol):
	def create(self, entry: Image) -> None: ...
