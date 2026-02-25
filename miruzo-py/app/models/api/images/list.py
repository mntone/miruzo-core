from collections.abc import Sequence
from typing import Annotated, final

from pydantic import ConfigDict, Field

from app.config.constants import INGEST_ID_MAXIMUM, INGEST_ID_MINIMUM
from app.models.api.variants.variant import VariantLayersModelBase, VariantModel
from app.models.records import ImageRecord
from app.models.types import VariantEntry


@final
class ImageListModel(VariantLayersModelBase):
	"""Summary metadata emitted by the list API."""

	model_config = ConfigDict(
		title='Image list model',
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
