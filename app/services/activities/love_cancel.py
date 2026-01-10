from datetime import datetime, time
from zoneinfo import ZoneInfo

from sqlmodel import Session

from app.config.score import ScoreConfig
from app.domain.activities.daily_period import resolve_daily_period_start
from app.domain.score.calculator import ScoreCalculator
from app.errors import InvalidStateError
from app.models.api.activities.responses import LoveStatsResponse
from app.models.api.activities.stats import LoveStatsModel
from app.models.enums import ActionKind
from app.persist.actions.factory import create_action_repository
from app.persist.actions.protocol import ActionRepository
from app.persist.stats.factory import create_stats_repository
from app.persist.users.factory import create_user_repository
from app.services.activities.actions.creator import ActionCreator
from app.services.activities.stats.score_factory import make_score_context


class LoveCancelRunner:
	def __init__(
		self,
		*,
		base_timezone: ZoneInfo | None,
		daily_reset_at: time,
		score_config: ScoreConfig,
	) -> None:
		self._base_timezone = base_timezone
		self._daily_reset_at = daily_reset_at
		self._score_calc = ScoreCalculator(score_config)

	def _resolve_last_loved_at(
		self,
		action_repo: ActionRepository,
		*,
		ingest_id: int,
		period_start: datetime,
	) -> datetime | None:
		pending_cancels = 0
		cutoff = period_start

		while True:
			action = action_repo.select_latest_one_by_multiple_kinds(
				ingest_id,
				kinds=(ActionKind.LOVE, ActionKind.LOVE_CANCELED),
				until_occurred_at=cutoff,
			)
			if action is None:
				return None

			if action.kind is ActionKind.LOVE_CANCELED:
				pending_cancels += 1
				cutoff = action.occurred_at
				continue

			if pending_cancels > 0:
				pending_cancels -= 1
				cutoff = action.occurred_at
				continue

			return action.occurred_at

	def run(self, session: Session, *, ingest_id: int, evaluated_at: datetime) -> LoveStatsResponse:
		# --- resolve period start ---
		period_start = resolve_daily_period_start(
			evaluated_at,
			daily_reset_at=self._daily_reset_at,
			base_timezone=self._base_timezone,
		)

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
		last_loved_at = self._resolve_last_loved_at(
			action_repo,
			ingest_id=ingest_id,
			period_start=period_start,
		)
		if last_loved_at is not None:
			stats.last_loved_at = last_loved_at

			if stats.first_loved_at > last_loved_at:
				stats.first_loved_at = last_loved_at
		else:
			stats.first_loved_at = None

		# --- update usage ---
		user_repo = create_user_repository(session)
		user = user_repo.get_or_create_singleton()
		if user.daily_love_used > 0:
			user.daily_love_used -= 1

		# --- create response ---
		response = LoveStatsResponse(
			stats=LoveStatsModel.from_record(stats),
		)

		return response
