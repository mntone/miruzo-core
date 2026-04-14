from datetime import datetime


class FixedClockProvider:
	def __init__(self, now: datetime) -> None:
		self._now = now

	def now(self) -> datetime:
		return self._now
