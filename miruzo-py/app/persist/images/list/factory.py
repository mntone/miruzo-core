from sqlmodel import Session

from app.persist.images.list.base import BaseImageListRepository
from app.persist.images.list.protocol import ImageListRepository


def create_image_list_repository(session: Session, *, engaged_score_threshold: int) -> ImageListRepository:
	"""
	Build an image list repository implementation for the configured backend.

	Args:
		session: SQLModel session bound to the current database engine.
		engaged_score_threshold: Minimum score_evaluated for engaged lists.

	Returns:
		Concrete repository tied to the active backend.
	"""

	return BaseImageListRepository(session, engaged_score_threshold=engaged_score_threshold)
