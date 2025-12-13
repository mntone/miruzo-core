from datetime import datetime
from typing import Set

from app.models.api.images.list import ImageListModel
from app.models.api.images.responses import ContextResponse, ImageListResponse
from app.services.images.repository.base import ImageRepository
from app.services.images.variants import compute_allowed_formats, normalize_variants_for_format


class ImageService:
	def __init__(self, repository: ImageRepository) -> None:
		self._repository = repository

	def get_latest(
		self,
		*,
		cursor: datetime | None,
		limit: int,
		exclude_formats: Set[str],
	) -> ImageListResponse:
		"""
		Return a paginated list of images (summary only).
		Includes variant normalization using allowed_formats.
		"""
		image_records, next_cursor = self._repository.get_list(cursor=cursor, limit=limit)

		allowed_formats = compute_allowed_formats(exclude_formats)

		output_images = []
		for image in image_records:
			normalized_layers = normalize_variants_for_format(image.variants, allowed_formats)
			image_model = ImageListModel.from_record(image, normalized_layers)
			output_images.append(image_model)

		response = ImageListResponse(
			items=output_images,
			cursor=next_cursor,
		)

		return response

	def get_context(
		self,
		image_id: int,
	) -> ContextResponse | None:
		"""
		Return a single image detail payload.

		Fetches the image record, increments view stats, and normalizes
		variant layers based on allowed formats.
		"""
		image = self._repository.get_detail(image_id)

		if image is None:
			return None

		stats = self._repository.upsert_stats_with_increment(image_id)

		response = ContextResponse.from_record(image, stats)

		return response
