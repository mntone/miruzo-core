from datetime import datetime

from app.domain.activities.daily_period import DailyPeriodResolver
from app.domain.decay_score.calculator import DecayScoreCalculator
from app.models.records import StatsRecord
from app.services.activities.stats.decay_score_factory import make_decay_score_context


def update_decay_score(
	*,
	stats: StatsRecord,
	evaluated_at: datetime,
	resolver: DailyPeriodResolver,
	decay_score_calculator: DecayScoreCalculator,
) -> int:
	context = make_decay_score_context(
		stats=stats,
		evaluated_at=evaluated_at,
		resolver=resolver,
	)

	new_score = decay_score_calculator.apply(
		score=stats.score,
		context=context,
	)

	stats.score = new_score
	stats.score_evaluated = new_score
	stats.score_evaluated_at = evaluated_at

	return new_score
