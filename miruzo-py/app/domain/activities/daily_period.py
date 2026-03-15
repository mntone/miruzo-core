from datetime import datetime, timedelta
from zoneinfo import ZoneInfo

_UTC_TIMEZONE = ZoneInfo('UTC')
_DEFAULT_DAY_START_OFFSET_LOCAL = timedelta(hours=5)


def _get_day_start_offset_utc(
	location: str,
) -> timedelta:
	# Use a fixed date so DST is intentionally ignored for the initial offset.
	timezone_offset = datetime(2026, 1, 1).replace(tzinfo=ZoneInfo(location)).utcoffset() or timedelta()
	day_start_offset = _DEFAULT_DAY_START_OFFSET_LOCAL - timezone_offset
	if day_start_offset.total_seconds() < 0:
		day_start_offset += timedelta(days=1)
	return day_start_offset


class DailyPeriodResolver:
	"""Resolve daily-period boundaries in UTC with a fixed day-start offset."""

	_day_start_offset: timedelta

	def __init__(self, day_start_offset: timedelta) -> None:
		self._day_start_offset = day_start_offset

	@property
	def day_start_offset(self) -> timedelta:
		return self._day_start_offset

	def resolve_period_start(self, evaluated_at: datetime) -> datetime:
		"""
		Return the UTC start datetime of the daily period that evaluated_at belongs to.
		Naive evaluated_at is treated as UTC.
		"""

		if evaluated_at.tzinfo is None:
			evaluated_at = evaluated_at.replace(tzinfo=_UTC_TIMEZONE)
		else:
			evaluated_at = evaluated_at.astimezone(_UTC_TIMEZONE)

		shifted_datetime = evaluated_at - self._day_start_offset
		shifted_date = shifted_datetime.replace(
			hour=0,
			minute=0,
			second=0,
			microsecond=0,
		)
		candidate = shifted_date + self._day_start_offset
		return candidate

	def resolve_period_end(self, evaluated_at: datetime) -> datetime:
		"""Return the UTC end datetime of the daily period for evaluated_at."""

		return self.resolve_period_start(evaluated_at) + timedelta(days=1)

	def resolve_period_range(self, evaluated_at: datetime) -> tuple[datetime, datetime]:
		"""
		Return the [since, until) datetime range of the daily period
		that evaluated_at belongs to in UTC.
		"""

		period_start = self.resolve_period_start(evaluated_at)

		period_end = period_start + timedelta(days=1)

		return period_start, period_end

	def is_since_period_start(
		self,
		target: datetime | None,
		*,
		evaluated_at: datetime,
	) -> bool:
		if target is None:
			return False

		period_start = self.resolve_period_start(evaluated_at)

		is_since_period = target >= period_start

		return is_since_period

	@classmethod
	def from_location(cls, location: str) -> 'DailyPeriodResolver':
		"""
		Build a resolver using a fixed-date UTC offset for the given location.
		The fixed date intentionally ignores DST.
		"""

		day_start_offset = _get_day_start_offset_utc(location)
		return cls(day_start_offset)
