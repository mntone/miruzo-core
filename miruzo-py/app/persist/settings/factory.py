from sqlmodel import Session

from app.config.environments import DatabaseBackend, env
from app.persist.settings.implementation import PostgreSQLSettingsRepository, SQLiteSettingsRepository
from app.persist.settings.protocol import SettingsRepository


def create_settings_repository(session: Session) -> SettingsRepository:
	"""
	Build a settings repository implementation for the configured backend.

	Args:
		session: SQLModel session bound to the current database engine.

	Returns:
		Concrete repository tied to the active backend.
	"""

	if env.database_backend == DatabaseBackend.SQLITE:
		return SQLiteSettingsRepository(session)
	elif env.database_backend == DatabaseBackend.POSTGRE_SQL:
		return PostgreSQLSettingsRepository(session)
	else:
		raise ValueError(f'Unsupported database type: {env.database_backend}')
