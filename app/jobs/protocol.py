from datetime import datetime
from typing import Protocol


class Job(Protocol):
	"""Executable job interface for scheduled tasks."""

	@property
	def name(self) -> str: ...

	def run(self, *, evaluated_at: datetime) -> None: ...
