from datetime import datetime
from typing import Annotated, final

from pydantic import BaseModel, ConfigDict, Field

from app.config.constants import SCORE_MAXIMUM, SCORE_MINIMUM
from app.models.records import StatsRecord


@final
class StatsModel(BaseModel):
	"""Aggregate engagement data for a single image."""

	model_config = ConfigDict(
		title='Image stats model',
		extra='forbid',
		frozen=True,
	)

	is_favorited: Annotated[
		bool,
		Field(
			title='Favorite flag',
			description='flag indicating whether the current user marked the image as a favorite',
		),
	] = False
	"""flag indicating whether the current user marked the image as a favorite"""

	score: Annotated[
		int,
		Field(
			title='Score',
			description=f'user-tunable ranking value clamped between {SCORE_MINIMUM} and {SCORE_MAXIMUM}',
			ge=SCORE_MINIMUM,
			le=SCORE_MAXIMUM,
		),
	]
	"""user-tunable ranking value"""

	view_count: Annotated[
		int,
		Field(title='View count', description='how many times this image has been viewed', ge=0),
	] = 0
	"""how many times this image has been viewed"""

	last_viewed_at: Annotated[
		datetime | None,
		Field(
			title='Last viewed timestamp',
			description="timestamp of the most recent view, or `null` if it hasn't been viewed yet",
		),
	] = None
	"""timestamp of the most recent view, or `None` if it hasn't been viewed yet"""

	@classmethod
	def from_record(cls, stats: StatsRecord) -> 'StatsModel':
		return cls(
			is_favorited=stats.favorite,
			score=stats.score,
			view_count=stats.view_count,
			last_viewed_at=stats.last_viewed_at,
		)
