from datetime import datetime
from typing import Annotated, final

from pydantic import BaseModel, ConfigDict, Field

from app.models.enums import ImageStatus
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
		Field(title='Image identifier', description='numeric primary key assigned in the database.'),
	]
	"""numeric primary key assigned in the database."""

	status: Annotated[
		ImageStatus,
		Field(title='Image status', description='lifecycle state: 0=active, 1=deleted, 2=missing'),
	] = ImageStatus.ACTIVE
	"""lifecycle state (active, deleted, missing) defined by `ImageStatus`"""

	captured_at: Annotated[
		datetime | None,
		Field(
			title='Captured timestamp',
			description='when the photo was originally shot; `null` if metadata is missing',
		),
	]
	"""when the photo was originally shot; `None` if metadata is missing"""

	ingested_at: Annotated[
		datetime,
		Field(title='Ingested timestamp', description='timestamp for when this record entered the system'),
	]
	"""timestamp for when this record entered the system"""

	@classmethod
	def from_record(cls, image: ImageRecord) -> 'SummaryModel':
		return cls(
			id=image.id,
			status=image.status,
			captured_at=image.captured_at,
			ingested_at=image.ingested_at,
		)
