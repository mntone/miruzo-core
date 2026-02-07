from datetime import datetime

from app.domain.activities.daily_period import DailyPeriodResolver
from app.domain.score.calculator import ScoreCalculator
from app.models.records import ActionRecord, StatsRecord
from app.services.activities.stats.score_factory import make_score_context


def update_score_from_action(
	*,
	stats: StatsRecord,
	action: ActionRecord,
	evaluated_at: datetime,
	resolver: DailyPeriodResolver,
	score_calculator: ScoreCalculator,
	update_evaluated: bool = False,
) -> int:
	context = make_score_context(
		stats=stats,
		evaluated_at=evaluated_at,
		resolver=resolver,
	)

	new_score = score_calculator.apply(
		action=action,
		score=stats.score,
		context=context,
	)

	stats.score = new_score
	if update_evaluated:
		stats.score_evaluated = new_score
		stats.score_evaluated_at = evaluated_at

	return new_score
