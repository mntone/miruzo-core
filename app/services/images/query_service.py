from collections.abc import Sequence
from datetime import datetime
from typing import cast, final

from sqlmodel import Session

import app.services.images.paginator as paginator
from app.config.variant import VariantLayerSpec
from app.models.api.images.responses import ImageListResponse
from app.models.records import ImageRecord
from app.persist.images.protocol import ImageRepository
from app.services.images.mapper import map_image_records_to_list_response
from app.services.images.query_executor import ImageListQueryExecutor


@final
class ImageQueryService:
	"""Query images and map records into API responses."""

	def __init__(
		self,
		*,
		session: Session,
		repository: ImageRepository,
		engaged_score_threshold: int,
		variant_layers: Sequence[VariantLayerSpec],
	) -> None:
		self._session = session
		self._repository = repository
		self._executor = ImageListQueryExecutor(session, engaged_score_threshold=engaged_score_threshold)
		self._variant_layers = variant_layers

	def get_latest(
		self,
		*,
		cursor: datetime | None,
		limit: int,
		exclude_formats: tuple[str, ...],
	) -> ImageListResponse[datetime]:
		"""
		Return a paginated list of images (summary only).
		Includes variant normalization using allowed_formats.
		"""

		# fmt: off
		rows = cast(
			Sequence[ImageRecord],
			self._executor
				.latest(cursor=cursor)
				.order_by_ingest_id()
				.limit(limit + 1)
				.execute(),
		)
		# fmt: on

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
		"""Return a paginated list of timeline images."""

		# fmt: off
		rows = cast(
			Sequence[tuple[ImageRecord, datetime]],
			self._executor
				.chronological(cursor=cursor)
				.order_by_ingest_id()
				.limit(limit + 1)
				.execute(),
		)
		# fmt: on

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
		"""Return a paginated list of recently viewed images."""

		# fmt: off
		rows = cast(
			Sequence[tuple[ImageRecord, datetime]],
			self._executor
				.recently(cursor=cursor)
				.order_by_ingest_id()
				.limit(limit + 1)
				.execute(),
		)
		# fmt: on

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
		"""Return a paginated list of first-loved images."""

		# fmt: off
		rows = cast(
			Sequence[tuple[ImageRecord, datetime]],
			self._executor
				.first_love(cursor=cursor)
				.order_by_ingest_id()
				.limit(limit + 1)
				.execute(),
		)
		# fmt: on

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
		"""Return a paginated list of hall-of-fame images."""

		# fmt: off
		rows = cast(
			Sequence[tuple[ImageRecord, datetime]],
			self._executor
				.hall_of_fame(cursor=cursor)
				.order_by_ingest_id()
				.limit(limit + 1)
				.execute(),
		)
		# fmt: on

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
		"""Return a paginated list of engaged images."""

		# fmt: off
		rows = cast(
			Sequence[tuple[ImageRecord, int]],
			self._executor
				.engaged(cursor=cursor)
				.order_by_ingest_id()
				.limit(limit + 1)
				.execute(),
		)
		# fmt: on

		items, next_cursor = paginator.slice_with_tuple_cursor(rows, limit)

		response = map_image_records_to_list_response(
			items,
			next_cursor=next_cursor,
			exclude_formats=exclude_formats,
			variant_layers=self._variant_layers,
		)

		return response
