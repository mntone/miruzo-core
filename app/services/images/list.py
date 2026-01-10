from collections.abc import Sequence
from datetime import datetime
from typing import final

import app.services.images.paginator as paginator
from app.config.variant import VariantLayerSpec
from app.models.api.images.responses import ImageListResponse
from app.persist.images.list.protocol import ImageListRepository
from app.services.images.mapper import map_image_records_to_list_response


@final
class ImageListService:
	"""Query images and map records into API responses."""

	def __init__(
		self,
		*,
		repository: ImageListRepository,
		variant_layers: Sequence[VariantLayerSpec],
	) -> None:
		self._repository = repository
		self._variant_layers = variant_layers

	def get_latest(
		self,
		*,
		cursor: datetime | None,
		limit: int,
		exclude_formats: tuple[str, ...],
	) -> ImageListResponse[datetime]:
		"""
		Return a paginated list of latest images.
		Includes variant normalization using allowed_formats.
		"""

		rows = self._repository.select_latest(cursor=cursor, limit=limit + 1)

		items, next_cursor = paginator.slice_with_cursor_latest(rows, limit)

		response = map_image_records_to_list_response(
			items,
			next_cursor=next_cursor,
			exclude_formats=exclude_formats,
			variant_layers=self._variant_layers,
		)

		return response

	def get_chronological(
		self,
		*,
		cursor: datetime | None,
		limit: int,
		exclude_formats: tuple[str, ...],
	) -> ImageListResponse[datetime]:
		"""
		Return a paginated list of timeline images.
		Includes variant normalization using allowed_formats.
		"""

		rows = self._repository.select_chronological(cursor=cursor, limit=limit + 1)

		items, next_cursor = paginator.slice_with_tuple_cursor(rows, limit)

		response = map_image_records_to_list_response(
			items,
			next_cursor=next_cursor,
			exclude_formats=exclude_formats,
			variant_layers=self._variant_layers,
		)

		return response

	def get_recently(
		self,
		*,
		cursor: datetime | None,
		limit: int,
		exclude_formats: tuple[str, ...],
	) -> ImageListResponse[datetime]:
		"""
		Return a paginated list of recently viewed images.
		Includes variant normalization using allowed_formats.
		"""

		rows = self._repository.select_recently(cursor=cursor, limit=limit + 1)

		items, next_cursor = paginator.slice_with_tuple_cursor(rows, limit)

		response = map_image_records_to_list_response(
			items,
			next_cursor=next_cursor,
			exclude_formats=exclude_formats,
			variant_layers=self._variant_layers,
		)

		return response

	def get_first_love(
		self,
		*,
		cursor: datetime | None,
		limit: int,
		exclude_formats: tuple[str, ...],
	) -> ImageListResponse[datetime]:
		"""
		Return a paginated list of first-loved images.
		Includes variant normalization using allowed_formats.
		"""

		rows = self._repository.select_first_love(cursor=cursor, limit=limit + 1)

		items, next_cursor = paginator.slice_with_tuple_cursor(rows, limit)

		response = map_image_records_to_list_response(
			items,
			next_cursor=next_cursor,
			exclude_formats=exclude_formats,
			variant_layers=self._variant_layers,
		)

		return response

	def get_hall_of_fame(
		self,
		*,
		cursor: datetime | None,
		limit: int,
		exclude_formats: tuple[str, ...],
	) -> ImageListResponse[datetime]:
		"""
		Return a paginated list of hall-of-fame images.
		Includes variant normalization using allowed_formats.
		"""

		rows = self._repository.select_hall_of_fame(cursor=cursor, limit=limit + 1)

		items, next_cursor = paginator.slice_with_tuple_cursor(rows, limit)

		response = map_image_records_to_list_response(
			items,
			next_cursor=next_cursor,
			exclude_formats=exclude_formats,
			variant_layers=self._variant_layers,
		)

		return response

	def get_engaged(
		self,
		*,
		cursor: int | None,
		limit: int,
		exclude_formats: tuple[str, ...],
	) -> ImageListResponse[int]:
		"""
		Return a paginated list of engaged images.
		Includes variant normalization using allowed_formats.
		"""

		rows = self._repository.select_engaged(cursor=cursor, limit=limit + 1)

		items, next_cursor = paginator.slice_with_tuple_cursor(rows, limit)

		response = map_image_records_to_list_response(
			items,
			next_cursor=next_cursor,
			exclude_formats=exclude_formats,
			variant_layers=self._variant_layers,
		)

		return response
