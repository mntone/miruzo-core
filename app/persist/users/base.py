# pyright: reportArgumentType=false

from abc import ABC, abstractmethod
from typing import TypeVar

from sqlalchemy import update
from sqlalchemy.exc import IntegrityError
from sqlmodel import Session, SQLModel

from app.models.records import UserRecord

TModel = TypeVar('TModel', bound=SQLModel)

_UNIQUE_USER_ID = 1


class BaseUserRepository(ABC):
	def __init__(self, session: Session) -> None:
		self._session = session

	@abstractmethod
	def _is_unique_violation(self, error: IntegrityError) -> bool: ...

	def _get_or_create(self, user_id: int) -> UserRecord:
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

	def get_or_create_singleton(self) -> UserRecord:
		user = self._get_or_create(_UNIQUE_USER_ID)

		return user

	def try_increment_daily_love_used(self, *, limit: int) -> bool:
		self._get_or_create(_UNIQUE_USER_ID)

		statement = (
			update(UserRecord)
			.where(
				UserRecord.id == _UNIQUE_USER_ID,
				UserRecord.daily_love_used < limit,
			)
			.values(daily_love_used=UserRecord.daily_love_used + 1)
		)

		result = self._session.exec(statement)

		return result.rowcount == 1

	def try_decrement_daily_love_used(self) -> bool:
		self._get_or_create(_UNIQUE_USER_ID)

		statement = (
			update(UserRecord)
			.where(
				UserRecord.id == _UNIQUE_USER_ID,
				UserRecord.daily_love_used > 0,
			)
			.values(daily_love_used=UserRecord.daily_love_used - 1)
		)

		result = self._session.exec(statement)

		return result.rowcount == 1

	def reset_daily_love_used(self) -> None:
		self._get_or_create(_UNIQUE_USER_ID)

		statement = update(UserRecord).where(UserRecord.id == _UNIQUE_USER_ID).values(daily_love_used=0)

		self._session.exec(statement)
