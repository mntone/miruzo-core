from datetime import datetime

from app.domain.activities.daily_period import DailyPeriodResolver
from app.domain.decay_score.context import DecayScoreContext
from app.errors import InvariantViolationError
from app.models.records import StatsRecord


def make_decay_score_context(
	*,
	stats: StatsRecord,
	evaluated_at: datetime,
	resolver: DailyPeriodResolver,
) -> DecayScoreContext:
	last_viewed_at = stats.last_viewed_at

	if last_viewed_at is None:
		days_since_last_view = 0
	elif last_viewed_at > evaluated_at:
		raise InvariantViolationError(
			f'last_viewed_at ({last_viewed_at}) is later than evaluated_at ({evaluated_at})',
		)
	else:
		days_since_last_view = (
			resolver.resolve_period_start(evaluated_at) - resolver.resolve_period_start(last_viewed_at)
		).days

	has_view_today = resolver.is_since_period_start(
		stats.last_viewed_at,
		evaluated_at=evaluated_at,
	)

	context = DecayScoreContext(
		days_since_last_view=days_since_last_view,
		has_view_today=has_view_today,
	)

	return context
