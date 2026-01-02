from datetime import datetime, time, timedelta, timezone

import pytest

from app.domain.activities.daily_period import (
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
	)

	assert result == datetime(2026, 1, 1, 5, 0, tzinfo=timezone.utc)


def test_resolve_daily_period_start_converts_timezone_to_utc() -> None:
	jst = timezone(timedelta(hours=9))
	evaluated_at = datetime(2026, 1, 2, 10, 0, tzinfo=jst)

	result = resolve_daily_period_start(
		evaluated_at,
		daily_reset_at=time(5, 0),
	)

	assert result == datetime(2026, 1, 1, 5, 0, tzinfo=timezone.utc)


def test_resolve_daily_period_start_assumes_utc_for_naive() -> None:
	evaluated_at = datetime(2026, 1, 2, 6, 0)

	result = resolve_daily_period_start(
		evaluated_at,
		daily_reset_at=time(5, 0),
	)

	assert result == datetime(2026, 1, 2, 5, 0, tzinfo=timezone.utc)


def test_resolve_daily_period_range_returns_one_day_span() -> None:
	evaluated_at = datetime(2026, 1, 2, 6, 0, tzinfo=timezone.utc)

	period_start, period_end = resolve_daily_period_range(
		evaluated_at,
		daily_reset_at=time(5, 0),
	)

	assert period_start == datetime(2026, 1, 2, 5, 0, tzinfo=timezone.utc)
	assert period_end == datetime(2026, 1, 3, 5, 0, tzinfo=timezone.utc)
