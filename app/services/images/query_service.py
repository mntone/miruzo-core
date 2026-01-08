from collections.abc import Sequence
from datetime import datetime
from typing import cast, final

from sqlmodel import Session

import app.services.images.paginator as paginator
from app.config.variant import VariantLayerSpec
from app.models.api.images.responses import ImageListResponse
from app.models.records import ImageRecord
from app.services.images.mapper import map_image_records_to_list_response
from app.services.images.query_executor import ImageListQueryExecutor
from app.services.images.repository import ImageRepository


@final
class ImageQueryService:
	"""Query images and map records into API responses."""

	def __init__(
		self,
		*,
		session: Session,
		repository: ImageRepository,
		variant_layers: Sequence[VariantLayerSpec],
	) -> None:
		self._session = session
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
		Return a paginated list of images (summary only).
		Includes variant normalization using allowed_formats.
		"""

		# fmt: off
		rows = cast(
			Sequence[ImageRecord],
			ImageListQueryExecutor(self._session)
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
			ImageListQueryExecutor(self._session)
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
			ImageListQueryExecutor(self._session)
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
			ImageListQueryExecutor(self._session)
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
			ImageListQueryExecutor(self._session)
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

	def get_by_ingest_id(self, ingest_id: int) -> ImageRecord | None:
		"""Return a single image record by ingest id."""

		record = self._repository.select_by_ingest_id(ingest_id)

		return record
