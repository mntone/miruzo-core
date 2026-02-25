from datetime import datetime
from typing import Annotated, final

from pydantic import BaseModel, ConfigDict, Field

from app.config.environments import env
from app.models.records import StatsRecord


@final
class LoveStatsModel(BaseModel):
	"""Love statistics for a single image."""

	model_config = ConfigDict(
		title='Love stats model',
		extra='forbid',
		frozen=True,
	)

	score: Annotated[
		int,
		Field(
			title='Score',
			description=f'user-tunable ranking value clamped between {env.score.minimum_score} and {env.score.maximum_score}',
			ge=env.score.minimum_score,
			le=env.score.maximum_score,
		),
	]
	"""user-tunable ranking value"""

	first_loved_at: Annotated[
		datetime | None,
		Field(
			title='First loved timestamp',
			description='timestamp of the first love action, or `null` if it has not been loved',
		),
	]
	"""timestamp of the first love action, or `None` if it has not been loved"""

	last_loved_at: Annotated[
		datetime | None,
		Field(
			title='Last loved timestamp',
			description='timestamp of the last love action, or `null` if it has not been loved',
		),
	]
	"""timestamp of the last love action, or `None` if it has not been loved"""

	@classmethod
	def from_record(cls, stats: StatsRecord) -> 'LoveStatsModel':
		return cls(
			score=stats.score,
			first_loved_at=stats.first_loved_at,
			last_loved_at=stats.last_loved_at,
		)


@final
class StatsModel(BaseModel):
	"""Aggregate engagement data for a single image."""

	model_config = ConfigDict(
		title='Stats model',
		extra='forbid',
		frozen=True,
	)

	score: Annotated[
		int,
		Field(
			title='Score',
			description=f'user-tunable ranking value clamped between {env.score.minimum_score} and {env.score.maximum_score}',
			ge=env.score.minimum_score,
			le=env.score.maximum_score,
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

	first_loved_at: Annotated[
		datetime | None,
		Field(
			title='First loved timestamp',
			description='timestamp of the first love action, or `null` if it has not been loved',
		),
	] = None
	"""timestamp of the first love action, or `None` if it has not been loved"""

	last_loved_at: Annotated[
		datetime | None,
		Field(
			title='Last loved timestamp',
			description='timestamp of the last love action, or `null` if it has not been loved',
		),
	] = None
	"""timestamp of the last love action, or `None` if it has not been loved"""

	hall_of_fame_at: Annotated[
		datetime | None,
		Field(
			title='Hall of fame timestamp',
			description='timestamp when the image entered the hall of fame, or `null` if it has not',
		),
	] = None
	"""timestamp when the image entered the hall of fame, or `None` if it has not"""

	view_milestone_count: Annotated[
		int | None,
		Field(
			title='View milestone',
			description='highest view milestone reached so far, or `null` if none',
			ge=1,
		),
	] = None
	"""highest view milestone reached so far, or `None` if none"""

	view_milestone_archived_at: Annotated[
		datetime | None,
		Field(
			title='View milestone timestamp',
			description='timestamp when the latest view milestone was reached, or `null` if none',
		),
	] = None
	"""timestamp when the latest view milestone was reached, or `None` if none"""

	@classmethod
	def from_record(cls, stats: StatsRecord) -> 'StatsModel':
		return cls(
			score=stats.score if stats.score >= env.score.minimum_score else env.score.minimum_score,
			view_count=stats.view_count,
			last_viewed_at=stats.last_viewed_at,
			first_loved_at=stats.first_loved_at,
			last_loved_at=stats.last_loved_at,
			hall_of_fame_at=stats.hall_of_fame_at,
			view_milestone_count=stats.view_milestone_count if stats.view_milestone_count != 0 else None,
			view_milestone_archived_at=stats.view_milestone_archived_at,
		)
