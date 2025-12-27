from datetime import datetime

from app.config.environments import env
from app.models.api.images.list import ImageListModel
from app.models.api.images.responses import ContextResponse, ImageListResponse
from app.services.images.repository.base import ImageRepository
from app.services.images.variants.api import compute_allowed_formats, normalize_variants_for_format
from app.services.images.variants.mapper import map_variants_to_layers


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
		image_records, next_cursor = self._repository.get_latest(cursor=cursor, limit=limit)

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

	def get_context(
		self,
		ingest_id: int,
	) -> ContextResponse | None:
		"""
		Return a single image detail payload.

		Fetches the image record, increments view stats, and normalizes
		variant layers based on allowed formats.
		"""
		image = self._repository.get_context(ingest_id)

		if image is None:
			return None

		stats = self._repository.upsert_stats_with_increment(ingest_id)

		response = ContextResponse.from_record(image, stats)

		return response
