from datetime import datetime, time, timezone
from zoneinfo import ZoneInfo

from tests.services.activities.stats.factory import build_stats_record

from app.services.activities.stats.score_factory import make_score_context


def test_make_score_context_without_last_view() -> None:
	stats = build_stats_record(1, last_viewed_at=None)
	evaluated_at = datetime(2024, 1, 2, 12, 0, tzinfo=timezone.utc)

	context = make_score_context(
		stats=stats,
		evaluated_at=evaluated_at,
		daily_reset_at=time(5, 0),
		base_timezone=ZoneInfo('UTC'),
	)

	assert context.last_viewed_at is None
	assert context.days_since_last_view == 0
	assert context.has_view_today is False


def test_make_score_context_marks_view_within_period() -> None:
	stats = build_stats_record(
		1,
		last_viewed_at=datetime(2024, 1, 1, 6, 0, tzinfo=timezone.utc),
	)
	evaluated_at = datetime(2024, 1, 2, 4, 0, tzinfo=timezone.utc)

	context = make_score_context(
		stats=stats,
		evaluated_at=evaluated_at,
		daily_reset_at=time(5, 0),
		base_timezone=ZoneInfo('UTC'),
	)

	assert context.has_view_today is True
	assert context.days_since_last_view == 0


def test_make_score_context_marks_view_before_period() -> None:
	stats = build_stats_record(
		1,
		last_viewed_at=datetime(2024, 1, 1, 4, 0, tzinfo=timezone.utc),
	)
	evaluated_at = datetime(2024, 1, 2, 4, 0, tzinfo=timezone.utc)

	context = make_score_context(
		stats=stats,
		evaluated_at=evaluated_at,
		daily_reset_at=time(5, 0),
		base_timezone=ZoneInfo('UTC'),
	)

	assert context.has_view_today is False
	assert context.days_since_last_view == 1
