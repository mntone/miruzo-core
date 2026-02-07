from datetime import datetime
from typing import final

from sqlmodel import Session

from app.domain.activities.daily_period import DailyPeriodResolver
from app.domain.score.calculator import ScoreCalculator
from app.persist.actions.factory import create_action_repository
from app.persist.stats.factory import create_stats_repository
from app.persist.users.factory import create_user_repository
from app.services.activities.actions.decay_creator import DecayActionCreator
from app.services.activities.stats.score_updater import update_score_from_action


@final
class DailyDecayRunner:
	"""Run daily decay for all stats using the provided dependencies."""

	def __init__(
		self,
		*,
		period_resolver: DailyPeriodResolver,
		score_calculator: ScoreCalculator,
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

			update_score_from_action(
				stats=stats,
				action=new_action,
				evaluated_at=evaluated_at,
				resolver=self._period_resolver,
				score_calculator=self._score_calculator,
				update_evaluated=True,
			)

		user_repo = create_user_repository(session)
		user = user_repo.get_or_create_singleton()
		user.daily_love_used = 0
