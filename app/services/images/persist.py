from app.models.records import ImageRecord
from app.services.images.repository.base import ImageRepository


class ImagePersistService:
	def __init__(self, repository: ImageRepository) -> None:
		self._repository = repository

	def record(self, image: ImageRecord) -> None:
		self._repository.insert(image)
