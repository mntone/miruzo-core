from abc import ABC, abstractmethod
from typing import final

from sqlalchemy.exc import IntegrityError
from sqlmodel import Session

from app.databases.mixins.postgre import PostgreSQLUniqueViolationMixin
from app.databases.mixins.sqlite import SQLiteUniqueViolationMixin
from app.models.records import SettingsRecord


class BaseSettingsRepository(ABC):
	def __init__(self, session: Session) -> None:
		self._session = session

	@abstractmethod
	def _is_unique_violation(self, error: IntegrityError) -> bool: ...

	def get(self, key: str) -> str | None:
		settings = self._session.get(SettingsRecord, key)
		if settings is None:
			return None

		return settings.value

	def insert(self, key: str, value: str) -> None:
		settings = SettingsRecord(key=key, value=value)
		self._session.add(settings)

		try:
			self._session.flush()
		except IntegrityError as exc:
			self._session.rollback()
			if not self._is_unique_violation(exc):
				raise


@final
class PostgreSQLSettingsRepository(PostgreSQLUniqueViolationMixin, BaseSettingsRepository): ...


@final
class SQLiteSettingsRepository(SQLiteUniqueViolationMixin, BaseSettingsRepository): ...
