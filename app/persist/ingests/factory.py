from sqlmodel import Session

from app.persist.ingests.base import BaseIngestRepository
from app.persist.ingests.protocol import IngestRepository


def create_ingest_repository(session: Session) -> IngestRepository:
	"""
	Build a ingest repository implementation for the configured backend.

	Args:
		session: SQLModel session bound to the current database engine.

	Returns:
		Concrete repository tied to the active backend.
	"""

	return BaseIngestRepository(session)
