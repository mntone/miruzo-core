from collections.abc import Sequence
from datetime import datetime

from app.domain.images.cursor import (
	DatetimeImageListCursor,
	ImageListCursorMode,
	UInt8ImageListCursor,
)
from app.models.records import ImageRecord


def slice_with_latest_cursor(
	rows: Sequence[ImageRecord],
	limit: int,
) -> tuple[Sequence[ImageRecord], DatetimeImageListCursor | None]:
	"""Trim to limit and return the next cursor for latest-ordered rows."""

	next_cursor: DatetimeImageListCursor | None = None
	if len(rows) > limit:
		last = rows[limit - 1]
		next_cursor = DatetimeImageListCursor(
			mode=ImageListCursorMode.LATEST,
			value=last.ingested_at,
			ingest_id=last.ingest_id,
		)

	items: Sequence[ImageRecord] = rows[:limit]

	return items, next_cursor


def slice_with_datetime_tuple_cursor(
	rows: Sequence[tuple[ImageRecord, datetime]],
	limit: int,
	*,
	mode: ImageListCursorMode,
) -> tuple[Sequence[ImageRecord], DatetimeImageListCursor | None]:
	"""Trim datetime-key rows and return the next cursor."""

	next_cursor: DatetimeImageListCursor | None = None
	if len(rows) > limit:
		image, value = rows[limit - 1]
		next_cursor = DatetimeImageListCursor(
			mode=mode,
			value=value,
			ingest_id=image.ingest_id,
		)

	items: list[ImageRecord] = [image for image, _ in rows[:limit]]

	return items, next_cursor


def slice_with_uint8_tuple_cursor(
	rows: Sequence[tuple[ImageRecord, int]],
	limit: int,
	*,
	mode: ImageListCursorMode,
) -> tuple[Sequence[ImageRecord], UInt8ImageListCursor | None]:
	"""Trim int-key rows and return the next cursor."""

	next_cursor: UInt8ImageListCursor | None = None
	if len(rows) > limit:
		image, value = rows[limit - 1]
		next_cursor = UInt8ImageListCursor(
			mode=mode,
			value=value,
			ingest_id=image.ingest_id,
		)

	items: list[ImageRecord] = [image for image, _ in rows[:limit]]

	return items, next_cursor
