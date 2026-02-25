from collections.abc import Sequence
from typing import Annotated, final

from pydantic import BaseModel, ConfigDict, Field

from app.models.api.activities.action import ActionModel
from app.models.api.activities.stats import StatsModel
from app.models.api.images.context import ImageRichModel, ImageSummaryModel


@final
class ContextResponse(BaseModel):
	"""Envelope returned by the context API for a single image."""

	model_config = ConfigDict(
		title='Image context response',
		extra='forbid',
		frozen=True,
	)

	image: Annotated[
		ImageSummaryModel | ImageRichModel,
		Field(
			title='Image summary or rich',
			description='metadata for the requested image',
		),
	]
	"""metadata for the requested image"""

	actions: Annotated[
		Sequence[ActionModel],
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

	@classmethod
	def from_record(
		cls,
		image: ImageSummaryModel | ImageRichModel,
		actions: Sequence[ActionModel],
		stats: StatsModel,
	) -> 'ContextResponse':
		return cls(
			image=image,
			actions=actions,
			stats=stats,
		)
