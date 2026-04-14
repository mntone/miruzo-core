from datetime import datetime, timezone

from app.domain.clock.protocol import ClockProvider


class _SystemClockProvider:
	def now(self) -> datetime:
		return datetime.now(timezone.utc)


def create_system_clock() -> ClockProvider:
	return _SystemClockProvider()
