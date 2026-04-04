from app.models.records import ImageRecord


class StubImageRepository:
	def __init__(self) -> None:
		self.one_response: ImageRecord | None = None
		self.one_called_with: int | None = None

		self.insert_called_with: ImageRecord | None = None

	def insert(self, image: ImageRecord) -> None:
		self.insert_called_with = image
