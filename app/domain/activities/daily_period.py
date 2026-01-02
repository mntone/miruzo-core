from datetime import datetime, time, timedelta, timezone


def resolve_daily_period_start(
	evaluated_at: datetime,
	*,
	daily_reset_at: time,
) -> datetime:
	"""
	Return the start datetime of the daily period
	that evaluated_at belongs to. evaluated_at is treated as UTC.
	If timezone-aware, it is converted to UTC.
	"""

	if evaluated_at.tzinfo is None:
		evaluated_at = evaluated_at.replace(tzinfo=timezone.utc)
	else:
		evaluated_at = evaluated_at.astimezone(timezone.utc)

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
) -> tuple[datetime, datetime]:
	"""
	Return the [since, until) datetime range of the daily period
	that evaluated_at belongs to.
	"""

	period_start = resolve_daily_period_start(
		evaluated_at,
		daily_reset_at=daily_reset_at,
	)

	period_end = period_start + timedelta(days=1)

	return period_start, period_end
