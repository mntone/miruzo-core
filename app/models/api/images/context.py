from collections.abc import Sequence
from datetime import datetime
from typing import Annotated, Literal, final

from pydantic import BaseModel, ConfigDict, Field

from app.config.constants import INGEST_ID_MAXIMUM, INGEST_ID_MINIMUM
from app.models.api.variants.variant import VariantLayersModelBase, VariantModel
from app.models.enums import ImageKind
from app.models.records import ImageRecord
from app.models.types import VariantEntry


class ImageSummaryModel(BaseModel):
	"""Lightweight metadata returned in context responses."""

	model_config = ConfigDict(
		title='Image summary model',
		extra='forbid',
		frozen=True,
	)

	id: Annotated[
		int,
		Field(
			title='Ingest identifier',
			description='numeric primary key assigned in the database.',
			ge=INGEST_ID_MINIMUM,
			le=INGEST_ID_MAXIMUM,
		),
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
	def from_record(cls, image: ImageRecord) -> 'ImageSummaryModel':
		return cls(
			id=image.ingest_id,
			ingested_at=image.ingested_at,
			kind=image.kind,
		)


@final
class ImageRichModel(ImageSummaryModel, VariantLayersModelBase):
	"""Rich metadata returned in context responses."""

	model_config = ConfigDict(
		title='Image rich model',
		extra='forbid',
		frozen=True,
	)

	level: Annotated[
		Literal['rich'],
		Field(
			title='Response level',
			description='response detail level for this context payload',
		),
	]
	"""response detail level for this context payload"""

	@classmethod
	def from_record(  # pyright: ignore[reportIncompatibleMethodOverride]
		cls,
		image: ImageRecord,
		normalized_layers: Sequence[Sequence[VariantEntry]],
	) -> 'ImageRichModel':
		# fmt: off
		return cls(
			level='rich',
			id=image.ingest_id,
			ingested_at=image.ingested_at,
			kind=image.kind,
			original=VariantModel.from_record(image.original),
			fallback=(
				VariantModel.from_record(image.fallback)
				if image.fallback is not None
				else None
			),
			variants=[
				[
					VariantModel.from_record(variant)
					for variant in layers
				]
				for layers in normalized_layers
			],
		)
		# fmt: on
