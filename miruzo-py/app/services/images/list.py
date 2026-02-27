from collections.abc import Sequence
from typing import final

import app.services.images.list_spec as spec
from app.config.variant import VariantLayerSpec
from app.domain.images.cursor import DatetimeImageListCursor, TImageListCursor, UInt8ImageListCursor
from app.models.api.images.responses import ImageListResponse
from app.persist.images.list.protocol import ImageListRepository
from app.services.images.list_spec import ImageListSpec, TRow
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

	def _get_list(
		self,
		list_spec: ImageListSpec[TRow, TImageListCursor],
		*,
		cursor: TImageListCursor | None,
		limit: int,
		exclude_formats: tuple[str, ...],
	) -> ImageListResponse[TImageListCursor]:
		rows = list_spec.fetch(self._repository, cursor, limit + 1)

		items, next_cursor = list_spec.slice(rows, limit)

		response = map_image_records_to_list_response(
			items,
			next_cursor=next_cursor,
			exclude_formats=exclude_formats,
			variant_layers=self._variant_layers,
		)

		return response

	def get_latest(
		self,
		*,
		cursor: DatetimeImageListCursor | None,
		limit: int,
		exclude_formats: tuple[str, ...],
	) -> ImageListResponse[DatetimeImageListCursor]:
		"""
		Return a paginated list of latest images.
		Includes variant normalization using allowed_formats.
		"""

		return self._get_list(
			spec.LATEST_SPEC,
			cursor=cursor,
			limit=limit,
			exclude_formats=exclude_formats,
		)

	def get_chronological(
		self,
		*,
		cursor: DatetimeImageListCursor | None,
		limit: int,
		exclude_formats: tuple[str, ...],
	) -> ImageListResponse[DatetimeImageListCursor]:
		"""
		Return a paginated list of timeline images.
		Includes variant normalization using allowed_formats.
		"""

		return self._get_list(
			spec.CHRONOLOGICAL_SPEC,
			cursor=cursor,
			limit=limit,
			exclude_formats=exclude_formats,
		)

	def get_recently(
		self,
		*,
		cursor: DatetimeImageListCursor | None,
		limit: int,
		exclude_formats: tuple[str, ...],
	) -> ImageListResponse[DatetimeImageListCursor]:
		"""
		Return a paginated list of recently viewed images.
		Includes variant normalization using allowed_formats.
		"""

		return self._get_list(
			spec.RECENTLY_SPEC,
			cursor=cursor,
			limit=limit,
			exclude_formats=exclude_formats,
		)

	def get_first_love(
		self,
		*,
		cursor: DatetimeImageListCursor | None,
		limit: int,
		exclude_formats: tuple[str, ...],
	) -> ImageListResponse[DatetimeImageListCursor]:
		"""
		Return a paginated list of first-loved images.
		Includes variant normalization using allowed_formats.
		"""

		return self._get_list(
			spec.FIRST_LOVE_SPEC,
			cursor=cursor,
			limit=limit,
			exclude_formats=exclude_formats,
		)

	def get_hall_of_fame(
		self,
		*,
		cursor: DatetimeImageListCursor | None,
		limit: int,
		exclude_formats: tuple[str, ...],
	) -> ImageListResponse[DatetimeImageListCursor]:
		"""
		Return a paginated list of hall-of-fame images.
		Includes variant normalization using allowed_formats.
		"""

		return self._get_list(
			spec.HALL_OF_FAME_SPEC,
			cursor=cursor,
			limit=limit,
			exclude_formats=exclude_formats,
		)

	def get_engaged(
		self,
		*,
		cursor: UInt8ImageListCursor | None,
		limit: int,
		exclude_formats: tuple[str, ...],
	) -> ImageListResponse[UInt8ImageListCursor]:
		"""
		Return a paginated list of engaged images.
		Includes variant normalization using allowed_formats.
		"""

		return self._get_list(
			spec.ENGAGED_SPEC,
			cursor=cursor,
			limit=limit,
			exclude_formats=exclude_formats,
		)
