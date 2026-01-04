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

	ingested_at: Annotated[
		datetime,
		Field(
			title='Ingested timestamp',
			description='timestamp when the image was ingested',
		),
	]
	"""timestamp when the image was ingested"""

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
			ingested_at=image.ingested_at,
			kind=image.kind,
		)
