from typing import Annotated, final

from pydantic import BaseModel, ConfigDict, Field

from app.config.constants import SCORE_MAXIMUM, SCORE_MINIMUM
from app.models.records import StatsRecord


@final
class FavoriteRequest(BaseModel):
	"""Payload for toggling an image favorite flag."""

	model_config = ConfigDict(
		title='Image favorite request',
		extra='forbid',
		strict=True,
	)

	value: Annotated[
		bool,
		Field(
			title='Favorite value',
			description='`true` to mark as favorite, `false` to remove the favorite flag',
		),
	]
	"""`true` to mark as favorite, `false` to remove the favorite flag"""


@final
class ScoreRequest(BaseModel):
	"""Payload for incrementing or decrementing an image score."""

	model_config = ConfigDict(
		title='Image score request',
		extra='forbid',
		strict=True,
	)

	delta: Annotated[
		int,
		Field(
			title='Score delta',
			description='Amount to add to the existing score (positive or negative)',
			ge=SCORE_MINIMUM - SCORE_MAXIMUM,
			le=SCORE_MAXIMUM - SCORE_MINIMUM,
		),
	]
	"""Amount to add to the existing score (positive or negative)"""


@final
class FavoriteResponse(BaseModel):
	"""Indicates the new favorite state after processing a request."""

	model_config = ConfigDict(
		title='Image favorite response',
		extra='forbid',
		frozen=True,
	)

	is_favorited: Annotated[
		bool,
		Field(
			title='Favorite state',
			description='`true` when the image is favorited after the update',
		),
	]
	"""`true` when the image is favorited after the update"""

	@classmethod
	def from_record(cls, stats: StatsRecord) -> 'FavoriteResponse':
		return cls(
			is_favorited=stats.favorite,
		)


@final
class ScoreResponse(BaseModel):
	"""Returns the updated score after applying a delta."""

	model_config = ConfigDict(
		title='Image score response',
		extra='forbid',
		frozen=True,
	)

	score: Annotated[
		int,
		Field(
			title='Current score',
			description='Resulting score clamped within the supported range',
			ge=SCORE_MINIMUM,
			le=SCORE_MAXIMUM,
		),
	]
	"""Resulting score clamped within the supported range"""

	@classmethod
	def from_record(cls, stats: StatsRecord) -> 'ScoreResponse':
		return cls(
			score=stats.score,
		)
