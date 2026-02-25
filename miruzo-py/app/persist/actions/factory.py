from sqlmodel import Session

from app.persist.actions.base import BaseActionRepository
from app.persist.actions.protocol import ActionRepository


def create_action_repository(session: Session) -> ActionRepository:
	"""
	Build a action repository implementation for the configured backend.

	Args:
		session: SQLModel session bound to the current database engine.

	Returns:
		Concrete repository tied to the active backend.
	"""

	return BaseActionRepository(session)
