from typing import final

from sqlmodel import Session

from app.models.records import ImageRecord


@final
class BaseImageRepository:
	def __init__(self, session: Session) -> None:
		self._session = session

	def select_by_ingest_id(
		self,
		ingest_id: int,
	) -> ImageRecord | None:
		image = self._session.get(ImageRecord, ingest_id)

		return image

	def insert(self, image: ImageRecord) -> None:
		self._session.add(image)
		self._session.flush()
		self._session.refresh(image)
