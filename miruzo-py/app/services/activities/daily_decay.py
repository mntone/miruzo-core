from datetime import datetime
from typing import final

from sqlmodel import Session

from app.domain.activities.daily_period import DailyPeriodResolver
from app.domain.decay_score.calculator import DecayScoreCalculator
from app.persist.actions.factory import create_action_repository
from app.persist.stats.factory import create_stats_repository
from app.persist.users.factory import create_user_repository
from app.services.activities.actions.decay_creator import DecayActionCreator
from app.services.activities.stats.decay_score_updater import update_decay_score


@final
class DailyDecayRunner:
	"""Run daily decay for all stats using the provided dependencies."""

	def __init__(
		self,
		*,
		period_resolver: DailyPeriodResolver,
		score_calculator: DecayScoreCalculator,
	) -> None:
		self._period_resolver = period_resolver
		self._score_calculator = score_calculator

	def apply_daily_decay(self, session: Session, *, evaluated_at: datetime) -> None:
		stats_repo = create_stats_repository(session)
		decay_creator = DecayActionCreator(
			repository=create_action_repository(session),
			period_resolver=self._period_resolver,
		)

		stats_list = stats_repo.iterable()
		for stats in stats_list:
			new_action = decay_creator.create(
				stats.ingest_id,
				occurred_at=evaluated_at,
			)

			# None means this decay has already been applied for the period
			if new_action is None:
				continue

			update_decay_score(
				stats=stats,
				evaluated_at=evaluated_at,
				resolver=self._period_resolver,
				decay_score_calculator=self._score_calculator,
			)

		user_repo = create_user_repository(session)
		user_repo.reset_daily_love_used()
