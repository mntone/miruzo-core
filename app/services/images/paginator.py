from collections.abc import Sequence
from datetime import datetime

from app.models.records import ImageRecord


def slice_with_cursor_latest(
	rows: Sequence[ImageRecord],
	limit: int,
) -> tuple[Sequence[ImageRecord], datetime | None]:
	"""Trim to limit and return the next cursor for latest-ordered rows."""

	next_cursor: datetime | None = None
	if len(rows) > limit:
		next_cursor = rows[limit - 1].ingested_at

	items: Sequence[ImageRecord] = rows[:limit]

	return items, next_cursor


def slice_with_cursor_for_datetime(
	rows: Sequence[tuple[ImageRecord, datetime]],
	limit: int,
) -> tuple[Sequence[ImageRecord], datetime | None]:
	"""Trim to limit and return the next cursor for datetime-ordered rows."""

	next_cursor: datetime | None = None
	if len(rows) > limit:
		_, next_cursor = rows[limit - 1]

	items: list[ImageRecord] = [image for image, _ in rows[:limit]]

	return items, next_cursor
