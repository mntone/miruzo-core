from collections.abc import Sequence
from typing import Annotated, final

from pydantic import BaseModel, ConfigDict, Field

from app.models.api.images.variant import VariantModel
from app.models.records import ImageRecord
from app.models.types import VariantEntry


@final
class ImageListModel(BaseModel):
	"""Summary metadata emitted by the list API."""

	model_config = ConfigDict(
		title='Image list model',
		extra='forbid',
		frozen=True,
	)

	id: Annotated[
		int,
		Field(title='Image identifier', description='numeric primary key assigned in the database.'),
	]
	"""numeric primary key assigned in the database."""

	original: Annotated[
		VariantModel,
		Field(
			title='Original variant',
			description='canonical full-resolution variant that all other renditions derive from',
		),
	]
	"""canonical full-resolution variant that all other renditions derive from"""

	fallback: Annotated[
		VariantModel | None,
		Field(
			title='Fallback variant',
			description="optional compatibility rendition used when layered variants can't be served",
		),
	] = None
	"""optional compatibility rendition used when layered variants can't be served"""

	variants: Annotated[
		Sequence[Sequence[VariantModel]],
		Field(
			title='Variant layers',
			description='layered list (e.g. primary/secondary) of alternative renditions organized by size',
		),
	]
	"""layered list (e.g. primary/secondary) of alternative renditions organized by size"""

	@classmethod
	def from_record(
		cls,
		image: ImageRecord,
		normalized_layers: Sequence[Sequence[VariantEntry]],
	) -> 'ImageListModel':
		# fmt: off
		return cls(
			id=image.ingest_id,

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
