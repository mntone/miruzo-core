from datetime import datetime
from typing import Annotated, final

from pydantic import BaseModel, ConfigDict, Field

from app.config.constants import LIMIT_MAXIMUM
from app.models.api.images.context.stats import StatsModel
from app.models.api.images.context.summary import SummaryModel
from app.models.api.images.list import ImageListModel
from app.models.records import ImageRecord, StatsRecord


@final
class ImageListResponse(BaseModel):
	"""Envelope returned by the latest-images API."""

	model_config = ConfigDict(
		title='Image list response',
		extra='forbid',
		frozen=True,
	)

	items: Annotated[
		list[ImageListModel],
		Field(
			title='Items',
			description='page of image summaries returned for this request',
			max_length=LIMIT_MAXIMUM,
		),
	]
	"""page of image summaries returned for this request"""

	cursor: Annotated[
		datetime | None,
		Field(
			title='Next cursor',
			default=None,
			description='pagination cursor to request the next page; `null` when no further pages exist',
		),
	]
	"""pagination cursor to request the next page; `None` when no further pages exist"""


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
		image: ImageRecord,
		stats: StatsRecord,
	) -> 'ContextResponse':
		return cls(
			image=SummaryModel.from_record(image),
			stats=StatsModel.from_record(stats),
		)
