from sqlmodel import Session

from app.core.settings import DatabaseBackend, settings
from app.services.images.repository.base import ImageRepository
from app.services.images.repository.postgre import PostgreSQLImageRepository
from app.services.images.repository.sqlite import SQLiteImageRepository


def create_image_repository(session: Session) -> ImageRepository:
	"""
	Build an image repository implementation for the configured backend.

	Args:
		session: SQLModel session bound to the current database engine.

	Returns:
		Concrete repository tied to the active backend.

	Raises:
		ValueError: if the configured backend is unsupported.
	"""
	if settings.database_backend == DatabaseBackend.SQLITE:
		return SQLiteImageRepository(session)
	elif settings.database_backend == DatabaseBackend.POSTGRE_SQL:
		return PostgreSQLImageRepository(session)
	else:
		raise ValueError(f'Unsupported database type: {settings.database_backend}')
