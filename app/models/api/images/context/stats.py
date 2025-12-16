from datetime import datetime
from typing import Annotated, Self, final

from pydantic import BaseModel, ConfigDict, Field

from app.config.constants import SCORE_MAXIMUM, SCORE_MINIMUM
from app.models.records import StatsRecord


@final
class StatsModel(BaseModel):
	model_config = ConfigDict(
		title='Image stats model',
		description='Aggregate engagement data for a single image.',
		validate_assignment=True,
	)

	is_favorited: Annotated[
		bool,
		Field(
			title='Favorite flag',
			description='flag indicating whether the current user marked the image as a favorite',
		),
	] = False
	score: Annotated[
		int,
		Field(
			title='Score',
			description=f'user-tunable ranking value clamped between {SCORE_MINIMUM} and {SCORE_MAXIMUM}',
			ge=SCORE_MINIMUM,
			le=SCORE_MAXIMUM,
		),
	]
	view_count: Annotated[
		int,
		Field(title='View count', description='how many times this image has been viewed', ge=0),
	] = 0
	last_viewed_at: Annotated[
		datetime | None,
		Field(
			title='Last viewed timestamp',
			description="timestamp of the most recent view, or `null` if it hasn't been viewed yet",
		),
	] = None

	@classmethod
	def from_record(cls, stats: StatsRecord) -> Self:
		return cls(
			is_favorited=stats.favorite,
			score=stats.score,
			view_count=stats.view_count,
			last_viewed_at=stats.last_viewed_at,
		)
