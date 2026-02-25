from sqlmodel import Session

from app.config.environments import DatabaseBackend, env
from app.persist.users.postgre import PostgreSQLUserRepository
from app.persist.users.protocol import UserRepository
from app.persist.users.sqlite import SQLiteUserRepository


def create_user_repository(session: Session) -> UserRepository:
	"""
	Build a user repository implementation for the configured backend.

	Args:
		session: SQLModel session bound to the current database engine.

	Returns:
		Concrete repository tied to the active backend.

	Raises:
		ValueError: if the configured backend is unsupported.
	"""

	if env.database_backend == DatabaseBackend.SQLITE:
		return SQLiteUserRepository(session)
	elif env.database_backend == DatabaseBackend.POSTGRE_SQL:
		return PostgreSQLUserRepository(session)
	else:
		raise ValueError(f'Unsupported database type: {env.database_backend}')
