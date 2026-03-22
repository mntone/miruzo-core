from datetime import datetime, timedelta, timezone

import pytest

from tests.services.activities.stats.factory import build_stats_record

from app.domain.activities.daily_period import DailyPeriodResolver
from app.errors import InvariantViolationError
from app.services.activities.stats.decay_score_factory import make_decay_score_context


def test_make_decay_score_contextwithout_last_view() -> None:
	stats = build_stats_record(1, last_viewed_at=None)
	evaluated_at = datetime(2024, 1, 2, 12, 0, tzinfo=timezone.utc)

	context = make_decay_score_context(
		stats=stats,
		evaluated_at=evaluated_at,
		resolver=DailyPeriodResolver(timedelta(hours=5)),
	)

	assert context.days_since_last_view == 0
	assert context.has_view_today is False


def test_make_decay_score_contextmarks_view_within_period() -> None:
	stats = build_stats_record(
		1,
		last_viewed_at=datetime(2024, 1, 1, 6, 0, tzinfo=timezone.utc),
	)
	evaluated_at = datetime(2024, 1, 2, 4, 0, tzinfo=timezone.utc)

	context = make_decay_score_context(
		stats=stats,
		evaluated_at=evaluated_at,
		resolver=DailyPeriodResolver(timedelta(hours=5)),
	)

	assert context.has_view_today is True
	assert context.days_since_last_view == 0


def test_make_decay_score_contextmarks_view_before_period() -> None:
	stats = build_stats_record(
		1,
		last_viewed_at=datetime(2024, 1, 1, 4, 0, tzinfo=timezone.utc),
	)
	evaluated_at = datetime(2024, 1, 2, 3, 0, tzinfo=timezone.utc)

	context = make_decay_score_context(
		stats=stats,
		evaluated_at=evaluated_at,
		resolver=DailyPeriodResolver(timedelta(hours=5)),
	)

	assert context.has_view_today is False
	assert context.days_since_last_view == 1


def test_make_decay_score_contextraises_for_future_last_view() -> None:
	stats = build_stats_record(1, last_viewed_at=datetime(2026, 1, 2, 12, 0, tzinfo=timezone.utc))
	evaluated_at = datetime(2026, 1, 2, 11, 0, tzinfo=timezone.utc)

	with pytest.raises(InvariantViolationError, match='last_viewed_at'):
		make_decay_score_context(
			stats=stats,
			evaluated_at=evaluated_at,
			resolver=DailyPeriodResolver(timedelta(hours=5)),
		)
