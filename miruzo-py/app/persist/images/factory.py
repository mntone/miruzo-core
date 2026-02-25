from sqlmodel import Session

from app.persist.images.base import BaseImageRepository
from app.persist.images.protocol import ImageRepository


def create_image_repository(session: Session) -> ImageRepository:
	"""
	Build a image repository implementation for the configured backend.

	Args:
		session: SQLModel session bound to the current database engine.

	Returns:
		Concrete repository tied to the active backend.
	"""

	return BaseImageRepository(session)
