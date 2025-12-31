from collections.abc import Sequence
from datetime import datetime
from typing import Annotated, final

from pydantic import BaseModel, ConfigDict, Field

from app.config.constants import LIMIT_MAXIMUM
from app.models.api.images.list import ImageListModel


@final
class ImageListResponse(BaseModel):
	"""Envelope returned by the latest-images API."""

	model_config = ConfigDict(
		title='Image list response',
		extra='forbid',
		frozen=True,
	)

	items: Annotated[
		Sequence[ImageListModel] | None,
		Field(
			title='Items',
			description='page of image summaries returned for this request',
			min_length=1,
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
