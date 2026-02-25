from datetime import datetime, time, timedelta, timezone
from zoneinfo import ZoneInfo, ZoneInfoNotFoundError

import pytest

from app.domain.activities.daily_period import (
	is_since_daily_period_start,
	resolve_daily_period_range,
	resolve_daily_period_start,
)


@pytest.mark.parametrize(
	'target_time',
	[
		time(6),
		time(5, 1),
		time(5, 0, 1),
		time(5, 0, 0, 1),
	],
)
def test_resolve_daily_period_start_keeps_time_after_reset(target_time: time) -> None:
	evaluated_at = datetime(2026, 1, 2, 6, 0, tzinfo=timezone.utc).replace(
		hour=target_time.hour,
		minute=target_time.minute,
		second=target_time.second,
		microsecond=target_time.microsecond,
	)

	result = resolve_daily_period_start(
		evaluated_at,
		daily_reset_at=time(5, 0),
		base_timezone=ZoneInfo('UTC'),
	)

	assert result == datetime(2026, 1, 2, 5, 0, tzinfo=timezone.utc)


@pytest.mark.parametrize(
	'target_time',
	[
		time(4),
		time(4, 59),
		time(4, 59, 59),
		time(4, 59, 59, 999_999),
	],
)
def test_resolve_daily_period_start_shifts_before_reset(target_time: time) -> None:
	evaluated_at = datetime(2026, 1, 2, 4, 0, tzinfo=timezone.utc).replace(
		hour=target_time.hour,
		minute=target_time.minute,
		second=target_time.second,
		microsecond=target_time.microsecond,
	)

	result = resolve_daily_period_start(
		evaluated_at,
		daily_reset_at=time(5, 0),
		base_timezone=ZoneInfo('UTC'),
	)

	assert result == datetime(2026, 1, 1, 5, 0, tzinfo=timezone.utc)


def test_resolve_daily_period_start_converts_timezone_to_utc() -> None:
	jst = timezone(timedelta(hours=9))
	evaluated_at = datetime(2026, 1, 2, 10, 0, tzinfo=jst)

	result = resolve_daily_period_start(
		evaluated_at,
		daily_reset_at=time(5, 0),
		base_timezone=ZoneInfo('UTC'),
	)

	assert result == datetime(2026, 1, 1, 5, 0, tzinfo=timezone.utc)


def test_resolve_daily_period_start_uses_base_timezone() -> None:
	try:
		jst = ZoneInfo('Asia/Tokyo')
	except ZoneInfoNotFoundError:
		pytest.skip('ZoneInfo data not available')
	evaluated_at = datetime(2026, 1, 2, 0, 30, tzinfo=timezone.utc)

	result = resolve_daily_period_start(
		evaluated_at,
		daily_reset_at=time(5, 0),
		base_timezone=jst,
	)

	assert result == datetime(2026, 1, 2, 5, 0, tzinfo=jst)


def test_resolve_daily_period_start_assumes_utc_for_naive() -> None:
	evaluated_at = datetime(2026, 1, 2, 6, 0)

	result = resolve_daily_period_start(
		evaluated_at,
		daily_reset_at=time(5, 0),
		base_timezone=ZoneInfo('UTC'),
	)

	assert result == datetime(2026, 1, 2, 5, 0, tzinfo=timezone.utc)


def test_resolve_daily_period_start_assumes_base_timezone_for_naive() -> None:
	try:
		jst = ZoneInfo('Asia/Tokyo')
	except ZoneInfoNotFoundError:
		pytest.skip('ZoneInfo data not available')
	evaluated_at = datetime(2026, 1, 2, 4, 0)

	result = resolve_daily_period_start(
		evaluated_at,
		daily_reset_at=time(5, 0),
		base_timezone=jst,
	)

	assert result == datetime(2026, 1, 1, 5, 0, tzinfo=jst)


def test_resolve_daily_period_start_handles_dst_transition() -> None:
	try:
		local_tz = ZoneInfo('America/New_York')
	except ZoneInfoNotFoundError:
		pytest.skip('ZoneInfo data not available')

	evaluated_at = datetime(2026, 3, 8, 10, 0, tzinfo=timezone.utc)

	result = resolve_daily_period_start(
		evaluated_at,
		daily_reset_at=time(5, 0),
		base_timezone=local_tz,
	)

	assert result == datetime(2026, 3, 8, 5, 0, tzinfo=local_tz)


def test_resolve_daily_period_range_returns_one_day_span() -> None:
	evaluated_at = datetime(2026, 1, 2, 6, 0, tzinfo=timezone.utc)

	period_start, period_end = resolve_daily_period_range(
		evaluated_at,
		daily_reset_at=time(5, 0),
		base_timezone=ZoneInfo('UTC'),
	)

	assert period_start == datetime(2026, 1, 2, 5, 0, tzinfo=timezone.utc)
	assert period_end == datetime(2026, 1, 3, 5, 0, tzinfo=timezone.utc)


def test_is_since_daily_period_start_handles_none() -> None:
	evaluated_at = datetime(2026, 1, 2, 6, 0, tzinfo=timezone.utc)

	result = is_since_daily_period_start(
		None,
		evaluated_at=evaluated_at,
		daily_reset_at=time(5, 0),
		base_timezone=ZoneInfo('UTC'),
	)
	assert result is False


@pytest.mark.parametrize(
	'target_time, expected',
	[
		(time(5, 0), True),
		(time(4, 59, 59), False),
		(time(5, 0, 1), True),
	],
)
def test_is_since_daily_period_start_checks_period_boundary(
	target_time: time,
	expected: bool,
) -> None:
	evaluated_at = datetime(2026, 1, 2, 6, 0, tzinfo=timezone.utc)
	target = datetime(2026, 1, 2, 0, 0, tzinfo=timezone.utc).replace(
		hour=target_time.hour,
		minute=target_time.minute,
		second=target_time.second,
		microsecond=target_time.microsecond,
	)

	result = is_since_daily_period_start(
		target,
		evaluated_at=evaluated_at,
		daily_reset_at=time(5, 0),
		base_timezone=ZoneInfo('UTC'),
	)
	assert result is expected
