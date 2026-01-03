from datetime import datetime, time
from zoneinfo import ZoneInfo

from app.domain.activities.daily_period import resolve_daily_period_start
from app.domain.score.context import ScoreContext
from app.errors import InvariantViolationError
from app.models.records import StatsRecord


def _has_view_in_current_period(
	*,
	last_viewed_at: datetime | None,
	evaluated_at: datetime,
	reset_time: time,
	timezone: ZoneInfo | None,
) -> bool:
	if last_viewed_at is None:
		return False

	period_start = resolve_daily_period_start(
		evaluated_at,
		daily_reset_at=reset_time,
		base_timezone=timezone,
	)

	return last_viewed_at >= period_start


def make_score_context(
	*,
	stats: StatsRecord,
	evaluated_at: datetime,
	daily_reset_at: time,
	base_timezone: ZoneInfo | None,
) -> ScoreContext:
	last_viewed_at = stats.last_viewed_at

	if last_viewed_at is None:
		days_since_last_view = 0
	elif last_viewed_at > evaluated_at:
		raise InvariantViolationError(
			f'last_viewed_at ({last_viewed_at}) is later than evaluated_at ({evaluated_at})',
		)
	else:
		days_since_last_view = (evaluated_at - last_viewed_at).days

	has_view_today = _has_view_in_current_period(
		last_viewed_at=stats.last_viewed_at,
		evaluated_at=evaluated_at,
		reset_time=daily_reset_at,
		timezone=base_timezone,
	)

	context = ScoreContext(
		evaluated_at=evaluated_at,
		last_viewed_at=last_viewed_at,
		days_since_last_view=days_since_last_view,
		has_view_today=has_view_today,
	)

	return context
