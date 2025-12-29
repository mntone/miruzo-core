from collections.abc import Sequence
from typing import Annotated, final

from pydantic import BaseModel, ConfigDict, Field

from app.models.api.activities.action import ActionModel
from app.models.api.activities.stats import StatsModel
from app.models.api.images.summary import SummaryModel
from app.models.records import ImageRecord, StatsRecord


@final
class ContextResponse(BaseModel):
	"""Envelope returned by the context API for a single image."""

	model_config = ConfigDict(
		title='Image context response',
		extra='forbid',
		frozen=True,
	)

	image: Annotated[
		SummaryModel,
		Field(
			title='Image summary',
			description='basic metadata for the requested image',
		),
	]
	"""basic metadata for the requested image"""

	actions: Annotated[
		Sequence[ActionModel] | None,
		Field(
			title='Actions',
			description='all actions',
			min_length=1,
		),
	]
	"""all actions"""

	stats: Annotated[
		StatsModel | None,
		Field(
			title='Image statistics',
			description='latest statistics for the image; `null` when stats are missing',
			default=None,
		),
	]
	"""latest statistics for the image; `None` when stats are missing"""

	# NOTE:
	# actions will become a required argument once context actions are wired.
	@classmethod
	def from_record(
		cls,
		image: ImageRecord,
		stats: StatsRecord,
	) -> 'ContextResponse':
		return cls(
			image=SummaryModel.from_record(image),
			actions=None,
			stats=StatsModel.from_record(stats),
		)
