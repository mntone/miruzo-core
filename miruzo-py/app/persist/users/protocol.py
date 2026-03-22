from typing import Protocol

from app.models.records import UserRecord


class UserRepository(Protocol):
	def create_singleton_if_missing(self) -> UserRecord: ...

	def get_singleton(self) -> UserRecord: ...

	def reset_daily_love_used(self) -> None: ...
