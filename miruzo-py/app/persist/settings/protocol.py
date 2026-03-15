from typing import Protocol


class SettingsRepository(Protocol):
	def get(self, key: str) -> str | None: ...

	def insert(self, key: str, value: str) -> None: ...
