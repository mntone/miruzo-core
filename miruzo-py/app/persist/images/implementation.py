from typing import final

from sqlalchemy import insert
from sqlalchemy.orm import Session

from app.databases.tables import image_table
from app.models.image import Image
from app.persist.images.protocol import ImageRepository


@final
class _ImageRepositoryImpl:
	def __init__(self, session: Session) -> None:
		self._session = session

	def create(self, entry: Image) -> None:
		stmt = insert(image_table).values(**entry.model_dump())
		self._session.execute(stmt)


def create_image_repository(session: Session) -> ImageRepository:
	"""
	Build an image repository implementation for the configured backend.

	Args:
		session: SQLAlchemy session bound to the current database engine.

	Returns:
		Concrete repository tied to the active backend.
	"""

	return _ImageRepositoryImpl(session)
