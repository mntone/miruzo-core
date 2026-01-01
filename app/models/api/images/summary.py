from datetime import datetime
from typing import Annotated, final

from pydantic import BaseModel, ConfigDict, Field

from app.models.enums import ImageKind
from app.models.records import ImageRecord


@final
class SummaryModel(BaseModel):
	"""Lightweight metadata returned in context responses."""

	model_config = ConfigDict(
		title='Image summary model',
		extra='forbid',
		frozen=True,
	)

	id: Annotated[
		int,
		Field(title='Ingest identifier', description='numeric primary key assigned in the database.'),
	]
	"""numeric primary key assigned in the database."""

	captured_at: Annotated[
		datetime,
		Field(
			title='Captured timestamp',
			description='timestamp when the photo was originally shot',
		),
	]
	"""timestamp when the photo was originally shot"""

	kind: Annotated[
		ImageKind,
		Field(
			title='Image kind',
			description='categorization of the image content (photo, illustration, or graphic)',
		),
	]
	"""categorization of the image content"""

	@classmethod
	def from_record(cls, image: ImageRecord) -> 'SummaryModel':
		return cls(
			id=image.ingest_id,
			captured_at=image.captured_at,
			kind=image.kind,
		)
