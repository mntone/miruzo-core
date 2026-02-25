from app.errors import SingletonUserMissingError
from app.models.records import UserRecord


class StubUserRepository:
	def __init__(self) -> None:
		self.users: dict[int, UserRecord] = {}
		self.get_called_with: list[int] = []

	def create_if_missing(self, user_id: int) -> UserRecord:
		self.get_called_with.append(user_id)

		user_record = self.users.get(user_id)
		if user_record is not None:
			return user_record

		user_record = UserRecord(id=user_id)
		self.users[user_id] = user_record
		return user_record

	def create_singleton_if_missing(self) -> UserRecord:
		return self.create_if_missing(1)

	def get_singleton(self) -> UserRecord:
		self.get_called_with.append(1)
		user_record = self.users.get(1)
		if user_record is None:
			raise SingletonUserMissingError('singleton user row is missing')

		return user_record

	def increment_daily_love_used(self, *, limit: int) -> bool:
		user_record = self.get_singleton()
		if user_record.daily_love_used >= limit:
			return False

		user_record.daily_love_used += 1
		return True

	def decrement_daily_love_used(self) -> bool:
		user_record = self.get_singleton()
		if user_record.daily_love_used <= 0:
			return False

		user_record.daily_love_used -= 1
		return True

	def reset_daily_love_used(self) -> None:
		user_record = self.get_singleton()
		user_record.daily_love_used = 0
