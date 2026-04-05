from app.models.image import Image


class StubImageRepository:
	def __init__(self) -> None:
		self.create_called_with: Image | None = None

	def create(self, entry: Image) -> None:
		self.create_called_with = entry
