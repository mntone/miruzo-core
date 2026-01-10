from datetime import datetime, time
from zoneinfo import ZoneInfo

from sqlmodel import Session

from app.config.score import ScoreConfig
from app.domain.activities.daily_period import resolve_daily_period_start
from app.domain.score.calculator import ScoreCalculator
from app.errors import InvalidStateError, QuotaExceededError
from app.models.api.activities.responses import LoveStatsResponse
from app.models.api.activities.stats import LoveStatsModel
from app.persist.actions.factory import create_action_repository
from app.persist.stats.factory import create_stats_repository
from app.persist.users.factory import create_user_repository
from app.services.activities.actions.creator import ActionCreator
from app.services.activities.stats.score_factory import make_score_context


class LoveRunner:
	def __init__(
		self,
		*,
		base_timezone: ZoneInfo | None,
		daily_reset_at: time,
		daily_love_limit: int,
		score_config: ScoreConfig,
	) -> None:
		self._base_timezone = base_timezone
		self._daily_reset_at = daily_reset_at
		self._daily_love_limit = daily_love_limit
		self._score_calc = ScoreCalculator(score_config)

	def run(self, session: Session, *, ingest_id: int, evaluated_at: datetime) -> LoveStatsResponse:
		# --- resolve period start ---
		period_start = resolve_daily_period_start(
			evaluated_at,
			daily_reset_at=self._daily_reset_at,
			base_timezone=self._base_timezone,
		)

		# --- try stats update ---
		stats_repo = create_stats_repository(session)
		updated = stats_repo.try_set_last_loved_at(
			ingest_id,
			last_loved_at=evaluated_at,
			since_occurred_at=period_start,
		)
		if not updated:
			raise InvalidStateError('already loved today')

		# --- load stats ---
		stats = stats_repo.get_one(ingest_id)

		# --- check quota ---
		user_repo = create_user_repository(session)
		user = user_repo.get_or_create_singleton()
		if user.daily_love_used >= self._daily_love_limit:
			raise QuotaExceededError()

		# --- insert action ---
		action_creator = ActionCreator(create_action_repository(session))
		new_action = action_creator.love(
			ingest_id,
			occurred_at=evaluated_at,
		)

		# --- update score ---
		context = make_score_context(
			stats=stats,
			evaluated_at=evaluated_at,
			daily_reset_at=self._daily_reset_at,
			base_timezone=self._base_timezone,
		)

		new_score = self._score_calc.apply(
			action=new_action,
			score=stats.score,
			context=context,
		)

		stats.score = new_score

		# --- update loved_at ---
		if stats.first_loved_at is None:
			stats.first_loved_at = evaluated_at
		stats.last_loved_at = evaluated_at

		# --- update usage ---
		user.daily_love_used += 1

		# --- create response ---
		response = LoveStatsResponse(
			stats=LoveStatsModel.from_record(stats),
		)

		return response
