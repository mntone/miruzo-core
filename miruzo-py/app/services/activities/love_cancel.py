from datetime import datetime

from sqlmodel import Session

from app.config.score import ScoreConfig
from app.domain.activities.daily_period import DailyPeriodResolver
from app.domain.score.calculator import ScoreCalculator
from app.errors import InvalidStateError
from app.models.api.activities.responses import LoveStatsResponse
from app.models.api.activities.stats import LoveStatsModel
from app.persist.actions.factory import create_action_repository
from app.persist.stats.factory import create_stats_repository
from app.persist.users.factory import create_user_repository
from app.services.activities.actions.creator import ActionCreator
from app.services.activities.stats.score_updater import update_score_from_action


class LoveCancelRunner:
	def __init__(
		self,
		*,
		period_resolver: DailyPeriodResolver,
		score_config: ScoreConfig,
	) -> None:
		self._period_resolver = period_resolver
		self._score_calc = ScoreCalculator(score_config)

	def run(self, session: Session, *, ingest_id: int, evaluated_at: datetime) -> LoveStatsResponse:
		# --- resolve period start ---
		period_start = self._period_resolver.resolve_period_start(evaluated_at)

		# --- try stats update ---
		stats_repo = create_stats_repository(session)
		updated = stats_repo.try_unset_last_loved_at(
			ingest_id,
			since_occurred_at=period_start,
		)
		if not updated:
			raise InvalidStateError('no love exists for today')

		# --- load stats ---
		stats = stats_repo.get_one(ingest_id)
		if stats.first_loved_at is None:
			raise InvalidStateError('first_loved_at is missing')

		# --- insert action ---
		action_repo = create_action_repository(session)
		action_creator = ActionCreator(action_repo)
		new_action = action_creator.cancel_love(
			ingest_id,
			occurred_at=evaluated_at,
		)

		# --- update score ---
		update_score_from_action(
			stats=stats,
			action=new_action,
			evaluated_at=evaluated_at,
			resolver=self._period_resolver,
			score_calculator=self._score_calc,
		)

		# --- update loved_at ---
		last_loved_action = action_repo.select_latest_effective_love(
			ingest_id,
			since_occurred_at=stats.first_loved_at,
			until_occurred_at=period_start,
		)
		if last_loved_action is not None:
			stats.last_loved_at = last_loved_action.occurred_at

			if stats.first_loved_at > last_loved_action.occurred_at:
				stats.first_loved_at = last_loved_action.occurred_at
		else:
			stats.first_loved_at = None

		# --- update usage ---
		user_repo = create_user_repository(session)
		user_repo.decrement_daily_love_used()

		# --- create response ---
		response = LoveStatsResponse(
			stats=LoveStatsModel.from_record(stats),
		)

		return response
