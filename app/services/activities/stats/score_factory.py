from datetime import datetime, time, timedelta

from app.domain.score.context import ScoreContext
from app.models.records import StatsRecord


def _compute_period_start(
	evaluated_at: datetime,
	reset: time,
) -> datetime:
	"""
	Return the start datetime of the daily scoring period
	that evaluated_at belongs to.
	"""

	candidate = evaluated_at.replace(
		hour=reset.hour,
		minute=reset.minute,
		second=reset.second,
		microsecond=0,
	)

	if evaluated_at < candidate:
		candidate -= timedelta(days=1)

	return candidate


def _has_view_in_current_period(
	*,
	last_viewed_at: datetime | None,
	evaluated_at: datetime,
	reset_time: time,
) -> bool:
	if last_viewed_at is None:
		return False

	period_start = _compute_period_start(evaluated_at, reset=reset_time)

	return last_viewed_at >= period_start


def make_score_context(
	*,
	stats: StatsRecord,
	evaluated_at: datetime,
	daily_reset_at: time,
) -> ScoreContext:
	last_viewed_at = stats.last_viewed_at

	if last_viewed_at is None:
		days_since_last_view = 0
	else:
		days_since_last_view = (evaluated_at - last_viewed_at).days

	has_view_today = _has_view_in_current_period(
		last_viewed_at=stats.last_viewed_at,
		evaluated_at=evaluated_at,
		reset_time=daily_reset_at,
	)

	context = ScoreContext(
		evaluated_at=evaluated_at,
		last_viewed_at=last_viewed_at,
		days_since_last_view=days_since_last_view,
		has_view_today=has_view_today,
	)

	return context
