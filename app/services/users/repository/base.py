from abc import ABC, abstractmethod
from typing import TypeVar

from sqlalchemy.exc import IntegrityError
from sqlmodel import Session, SQLModel

from app.models.records import UserRecord

TModel = TypeVar('TModel', bound=SQLModel)


class BaseUserRepository(ABC):
	def __init__(self, session: Session) -> None:
		self._session = session

	@abstractmethod
	def _is_unique_violation(self, error: IntegrityError) -> bool: ...

	def get_or_create(self, user_id: int) -> UserRecord:
		user = self._session.get(UserRecord, user_id)
		if user is not None:
			return user

		user = UserRecord(id=user_id)
		self._session.add(user)

		try:
			self._session.flush()
		except IntegrityError as exc:
			self._session.rollback()
			if not self._is_unique_violation(exc):
				raise
			user = self._session.get_one(UserRecord, user_id)

		return user
