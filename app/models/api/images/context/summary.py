from datetime import datetime
from typing import Annotated, Self, final

from pydantic import BaseModel, ConfigDict, Field

from app.models.enums import ImageStatus
from app.models.records import ImageRecord


@final
class SummaryModel(BaseModel):
	model_config = ConfigDict(
		title='Image summary model',
		description='Lightweight metadata returned in context responses.',
		validate_assignment=True,
	)

	id: Annotated[
		int,
		Field(title='Image identifier', description='numeric primary key assigned in the database.'),
	]
	status: Annotated[
		ImageStatus,
		Field(title='Image status', description='lifecycle state: 0=active, 1=deleted, 2=missing'),
	] = ImageStatus.ACTIVE

	captured_at: Annotated[
		datetime | None,
		Field(
			title='Captured timestamp',
			description='when the photo was originally shot; `null` if metadata is missing',
		),
	]
	ingested_at: Annotated[
		datetime,
		Field(title='Ingested timestamp', description='timestamp for when this record entered the system'),
	]

	@property
	def ingested_at(self) -> datetime:
		return self.ingested_at

	@classmethod
	def from_record(cls, image: ImageRecord) -> Self:
		return cls(
			id=image.id,
			status=image.status,
			captured_at=image.captured_at,
			ingested_at=image.ingested_at,
		)
