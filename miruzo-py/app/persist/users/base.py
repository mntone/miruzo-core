# pyright: reportAttributeAccessIssue=false
# pyright: reportArgumentType=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false
# pyright: reportUnknownVariableType=false

from abc import ABC, abstractmethod
from typing import TypeVar

from sqlalchemy import update
from sqlalchemy.exc import IntegrityError
from sqlmodel import Session, SQLModel

from app.errors import SingletonUserMissingError
from app.models.records import UserRecord

TModel = TypeVar('TModel', bound=SQLModel)

_UNIQUE_USER_ID = 1


class BaseUserRepository(ABC):
	def __init__(self, session: Session) -> None:
		self._session = session

	@abstractmethod
	def _is_unique_violation(self, error: IntegrityError) -> bool: ...

	def _expire_singleton_daily_love_used(self) -> None:
		# Core UPDATE bypasses ORM state sync, so expire only the touched field.
		for entity in self._session.identity_map.values():
			if isinstance(entity, UserRecord) and entity.id == _UNIQUE_USER_ID:
				self._session.expire(entity, ['daily_love_used'])
				return

	def _create_if_missing(self, user_id: int) -> UserRecord:
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

	def create_singleton_if_missing(self) -> UserRecord:
		return self._create_if_missing(_UNIQUE_USER_ID)

	def get_singleton(self) -> UserRecord:
		user = self._session.get(UserRecord, _UNIQUE_USER_ID)
		if user is None:
			raise SingletonUserMissingError('singleton user row is missing')

		return user

	def increment_daily_love_used(self, *, limit: int) -> bool:
		self._session.flush()

		user_table = UserRecord.__table__

		statement = (
			update(user_table)
			.where(
				user_table.c.id == _UNIQUE_USER_ID,
				user_table.c.daily_love_used < limit,
			)
			.values(daily_love_used=user_table.c.daily_love_used + 1)
			.returning(user_table.c.id)
		)

		updated_id = self._session.exec(statement).scalar_one_or_none()
		if updated_id is not None:
			self._expire_singleton_daily_love_used()
			return True

		# Distinguish quota exhaustion from missing singleton.
		self.get_singleton()

		return False

	def decrement_daily_love_used(self) -> bool:
		self._session.flush()

		user_table = UserRecord.__table__

		statement = (
			update(user_table)
			.where(
				user_table.c.id == _UNIQUE_USER_ID,
				user_table.c.daily_love_used > 0,
			)
			.values(daily_love_used=user_table.c.daily_love_used - 1)
			.returning(user_table.c.id)
		)

		updated_id = self._session.exec(statement).scalar_one_or_none()
		if updated_id is not None:
			self._expire_singleton_daily_love_used()
			return True

		# Distinguish lower-bound exhaustion from missing singleton.
		self.get_singleton()

		return False

	def reset_daily_love_used(self) -> None:
		self._session.flush()

		user_table = UserRecord.__table__

		statement = (
			update(user_table)
			.where(user_table.c.id == _UNIQUE_USER_ID)
			.values(daily_love_used=0)
			.returning(user_table.c.id)
		)

		updated_id = self._session.exec(statement).scalar_one_or_none()
		if updated_id is None:
			raise SingletonUserMissingError('singleton user row is missing')

		self._expire_singleton_daily_love_used()
