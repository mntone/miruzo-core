from typing import Annotated, Self, final

from pydantic import BaseModel, ConfigDict, Field

from app.models.api.images.variant import VariantModel
from app.models.enums import ImageStatus
from app.models.records import ImageRecord, VariantRecord


@final
class ImageListModel(BaseModel):
	model_config = ConfigDict(
		title='Image list model',
		description='Summary metadata emitted by the list API.',
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

	original: Annotated[
		VariantModel,
		Field(
			title='Original variant',
			description='canonical full-resolution variant that all other renditions derive from',
		),
	]
	fallback: Annotated[
		VariantModel | None,
		Field(
			title='Fallback variant',
			description="optional compatibility rendition used when layered variants can't be served",
		),
	] = None
	variants: Annotated[
		list[list[VariantModel]],
		Field(
			title='Variant layers',
			description='layered list (e.g. primary/secondary) of alternative renditions organized by size',
		),
	]

	@property
	def id(self) -> int:
		return self._id

	@classmethod
	def from_record(
cls,
image: ImageRecord,
normalized_layers: list[list[VariantRecord]],
) -> 'ImageListModel':
		# fmt: off
		return cls(
			id=image.id,
			status=image.status,

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
