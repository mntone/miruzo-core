from datetime import datetime
from typing import final

from app.config.environments import env
from app.models.api.images.list import ImageListModel
from app.models.api.images.responses import ImageListResponse
from app.models.records import ImageRecord
from app.services.images.repository import ImageRepository
from app.services.images.variants.api import compute_allowed_formats, normalize_variants_for_format
from app.services.images.variants.mapper import map_variants_to_layers


@final
class ImageQueryService:
	def __init__(self, repository: ImageRepository) -> None:
		self._repository = repository

	def get_latest(
		self,
		*,
		cursor: datetime | None,
		limit: int,
		exclude_formats: tuple[str, ...],
	) -> ImageListResponse:
		"""
		Return a paginated list of images (summary only).
		Includes variant normalization using allowed_formats.
		"""
		image_records, next_cursor = self._repository.select_latest(cursor=cursor, limit=limit)

		allowed_formats = compute_allowed_formats(exclude_formats)

		output_images: list[ImageListModel] = []
		for image in image_records:
			layers = map_variants_to_layers(image.variants, spec=env.variant_layers)
			normalized_layers = normalize_variants_for_format(layers, allowed_formats)
			image_model = ImageListModel.from_record(image, normalized_layers)
			output_images.append(image_model)

		response = ImageListResponse(
			items=output_images,
			cursor=next_cursor,
		)

		return response

	def get_by_ingest_id(self, ingest_id: int) -> ImageRecord | None:
		record = self._repository.select_by_ingest_id(ingest_id)

		return record
