from app.models.records import UserRecord


class StubUserRepository:
	def __init__(self) -> None:
		self.users: dict[int, UserRecord] = {}
		self.get_called_with: list[int] = []

	def get_or_create(self, user_id: int) -> UserRecord:
		self.get_called_with.append(user_id)

		user_record = self.users.get(user_id)
		if user_record is not None:
			return user_record

		user_record = UserRecord(id=user_id)
		self.users[user_id] = user_record
		return user_record
