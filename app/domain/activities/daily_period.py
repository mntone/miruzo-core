from datetime import datetime, time, timedelta
from typing import cast
from zoneinfo import ZoneInfo

_LOCAL_TIMEZONE = cast(ZoneInfo, datetime.now().astimezone().tzinfo or ZoneInfo('UTC'))


def resolve_daily_period_start(
	evaluated_at: datetime,
	*,
	daily_reset_at: time,
	base_timezone: ZoneInfo | None,
) -> datetime:
	"""
	Return the start datetime of the daily period that evaluated_at belongs to.
	evaluated_at is interpreted in base_timezone (local timezone if None).
	"""

	if base_timezone is None:
		base_timezone = _LOCAL_TIMEZONE

	if evaluated_at.tzinfo is None:
		evaluated_at = evaluated_at.replace(tzinfo=base_timezone)
	else:
		evaluated_at = evaluated_at.astimezone(base_timezone)

	candidate = evaluated_at.replace(
		hour=daily_reset_at.hour,
		minute=daily_reset_at.minute,
		second=daily_reset_at.second,
		microsecond=0,
	)

	if candidate > evaluated_at:
		candidate -= timedelta(days=1)

	return candidate


def resolve_daily_period_range(
	evaluated_at: datetime,
	*,
	daily_reset_at: time,
	base_timezone: ZoneInfo | None,
) -> tuple[datetime, datetime]:
	"""
	Return the [since, until) datetime range of the daily period
	that evaluated_at belongs to, using base_timezone or local timezone.
	"""

	period_start = resolve_daily_period_start(
		evaluated_at,
		daily_reset_at=daily_reset_at,
		base_timezone=base_timezone,
	)

	period_end = period_start + timedelta(days=1)

	return period_start, period_end


def is_since_daily_period_start(
	target: datetime | None,
	*,
	evaluated_at: datetime,
	daily_reset_at: time,
	base_timezone: ZoneInfo | None,
) -> bool:
	if target is None:
		return False

	period_start = resolve_daily_period_start(
		evaluated_at,
		daily_reset_at=daily_reset_at,
		base_timezone=base_timezone,
	)

	is_since_period = target >= period_start

	return is_since_period
