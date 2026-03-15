from typing import final


@final
class StubSettingsRepository:
	def __init__(self) -> None:
		self.store: dict[str, str] = {}
		self.inserts: list[tuple[str, str]] = []

	def get(self, key: str) -> str | None:
		return self.store.get(key)

	def insert(self, key: str, value: str) -> None:
		self.store[key] = value
		self.inserts.append((key, value))
