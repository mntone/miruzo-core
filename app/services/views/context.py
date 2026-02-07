from datetime import datetime, timezone
from typing import final

from sqlmodel import Session

from app.config.constants import VIEW_MILESTONES
from app.config.environments import Settings
from app.domain.activities.daily_period import DailyPeriodResolver
from app.domain.score.calculator import ScoreCalculator
from app.models.api.activities.action import ActionModel
from app.models.api.activities.stats import StatsModel
from app.models.api.context.query import ContextQuery
from app.models.api.context.responses import ContextResponse
from app.models.api.images.context import ImageRichModel, ImageSummaryModel
from app.persist.actions.protocol import ActionRepository
from app.persist.images.protocol import ImageRepository
from app.persist.stats.protocol import StatsRepository
from app.services.activities.actions.creator import ActionCreator
from app.services.activities.stats.score_updater import update_score_from_action
from app.services.images.variants.api import compute_allowed_formats, normalize_variants_for_format
from app.services.images.variants.mapper import map_variants_to_layers


@final
class ContextService:
	def __init__(
		self,
		session: Session,
		*,
		action_repo: ActionRepository,
		image_repo: ImageRepository,
		stats_repo: StatsRepository,
		env: Settings,
	) -> None:
		self._session = session
		self._action_repo = action_repo
		self._action_persist = ActionCreator(action_repo)
		self._image_repo = image_repo
		self._stats_repo = stats_repo
		self._score_calc = ScoreCalculator(env.score)
		self._period_resolver = DailyPeriodResolver(
			base_timezone=env.base_timezone,
			daily_reset_at=env.time.daily_reset_at,
		)
		self._variant_layers = env.variant_layers

	def get_context(
		self,
		ingest_id: int,
		*,
		query: ContextQuery,
	) -> ContextResponse | None:
		"""
		Return a single image detail payload.

		Fetches the image record, increments view stats, and normalizes
		variant layers based on allowed formats.
		"""

		with self._session.begin():
			image = self._image_repo.select_by_ingest_id(ingest_id)
			if image is None:
				return None

			current = datetime.now(timezone.utc)
			new_action = self._action_persist.view(
				ingest_id,
				occurred_at=current,
			)

			stats = self._stats_repo.get_or_create(
				ingest_id,
				initial_score=self._score_calc.config.initial_score,
			)

			update_score_from_action(
				stats=stats,
				action=new_action,
				evaluated_at=current,
				resolver=self._period_resolver,
				score_calculator=self._score_calc,
			)
			stats.view_count += 1
			stats.last_viewed_at = current

			for milestone in VIEW_MILESTONES:
				if stats.view_count >= milestone:
					if stats.view_milestone_count < milestone:
						stats.view_milestone_count = milestone
						stats.view_milestone_archived_at = current

		actions = self._action_repo.select_by_ingest_id(ingest_id)

		image_response: ImageSummaryModel | ImageRichModel
		match query.level:
			case 'default':
				image_response = ImageSummaryModel.from_record(image)
			case 'rich':
				layers = map_variants_to_layers(image.variants, spec=self._variant_layers)
				allowed_formats = compute_allowed_formats(query.exclude_formats)
				normalized_layers = normalize_variants_for_format(layers, allowed_formats)
				image_response = ImageRichModel.from_record(image, normalized_layers)

		response = ContextResponse.from_record(
			image=image_response,
			actions=[ActionModel.from_record(action) for action in actions],
			stats=StatsModel.from_record(stats),
		)

		return response
