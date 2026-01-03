from datetime import datetime, timezone
from typing import final

from sqlmodel import Session

from app.config.constants import VIEW_MILESTONES
from app.config.environments import Settings
from app.domain.score.calculator import ScoreCalculator
from app.models.api.activities.action import ActionModel
from app.models.api.activities.stats import StatsModel
from app.models.api.context.responses import ContextResponse
from app.models.api.images.summary import SummaryModel
from app.services.activities.actions.creator import ActionCreator
from app.services.activities.actions.repository import ActionRepository
from app.services.activities.stats.repository.protocol import StatsRepository
from app.services.activities.stats.score_factory import make_score_context
from app.services.images.query_service import ImageQueryService


@final
class ContextService:
	def __init__(
		self,
		session: Session,
		*,
		action: ActionRepository,
		image_query: ImageQueryService,
		stats: StatsRepository,
		env: Settings,
	) -> None:
		self._session = session
		self._action = action
		self._action_persist = ActionCreator(action)
		self._image = image_query
		self._stats = stats
		self._score_calc = ScoreCalculator(env.score)
		self._daily_reset_at = env.time.daily_reset_at
		self._base_timezone = env.base_timezone

	def get_context(
		self,
		ingest_id: int,
	) -> ContextResponse | None:
		"""
		Return a single image detail payload.

		Fetches the image record, increments view stats, and normalizes
		variant layers based on allowed formats.
		"""

		with self._session.begin():
			image = self._image.get_by_ingest_id(ingest_id)
			if image is None:
				return None

			current = datetime.now(timezone.utc)
			new_action = self._action_persist.view(
				ingest_id,
				occurred_at=current,
			)

			stats = self._stats.get_or_create(
				ingest_id,
				initial_score=self._score_calc.config.initial_score,
			)

			context = make_score_context(
				stats=stats,
				evaluated_at=current,
				daily_reset_at=self._daily_reset_at,
				base_timezone=self._base_timezone,
			)

			new_score = self._score_calc.apply(
				action=new_action,
				score=stats.score,
				context=context,
			)

			stats.score = new_score
			stats.view_count += 1
			stats.last_viewed_at = current

			for milestone in VIEW_MILESTONES:
				if stats.view_count >= milestone:
					if stats.view_milestone_count < milestone:
						stats.view_milestone_count = milestone
						stats.view_milestone_archived_at = current

		actions = self._action.select_by_ingest_id(ingest_id)

		response = ContextResponse.from_record(
			image=SummaryModel.from_record(image),
			actions=[ActionModel.from_record(action) for action in actions],
			stats=StatsModel.from_record(stats),
		)

		return response
