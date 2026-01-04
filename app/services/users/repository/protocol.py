from typing import Protocol

from app.models.records import UserRecord


class UserRepository(Protocol):
	def get_or_create(self, user_id: int) -> UserRecord: ...
