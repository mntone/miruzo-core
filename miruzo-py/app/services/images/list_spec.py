from collections.abc import Sequence
from dataclasses import dataclass
from datetime import datetime
from typing import Callable, Generic, TypeVar, final

import app.services.images.paginator as paginator
from app.domain.images.cursor import (
	DatetimeImageListCursor,
	ImageListCursorMode,
	TImageListCursor,
	UInt8ImageListCursor,
)
from app.models.records import ImageRecord
from app.persist.images.list.protocol import ImageListRepository

TRow = TypeVar('TRow')


@dataclass(frozen=True, slots=True)
@final
class ImageListSpec(Generic[TRow, TImageListCursor]):
	fetch: Callable[[ImageListRepository, TImageListCursor | None, int], Sequence[TRow]]
	slice: Callable[[Sequence[TRow], int], tuple[Sequence[ImageRecord], TImageListCursor | None]]


LATEST_SPEC: ImageListSpec[ImageRecord, DatetimeImageListCursor] = ImageListSpec(
	fetch=lambda repo, cursor, limit: repo.select_latest(cursor=cursor, limit=limit),
	slice=paginator.slice_with_latest_cursor,
)

CHRONOLOGICAL_SPEC: ImageListSpec[tuple[ImageRecord, datetime], DatetimeImageListCursor] = ImageListSpec(
	fetch=lambda repo, cursor, limit: repo.select_chronological(cursor=cursor, limit=limit),
	slice=lambda rows, limit: paginator.slice_with_datetime_tuple_cursor(
		rows,
		limit,
		mode=ImageListCursorMode.CHRONOLOGICAL,
	),
)

RECENTLY_SPEC: ImageListSpec[tuple[ImageRecord, datetime], DatetimeImageListCursor] = ImageListSpec(
	fetch=lambda repo, cursor, limit: repo.select_recently(cursor=cursor, limit=limit),
	slice=lambda rows, limit: paginator.slice_with_datetime_tuple_cursor(
		rows,
		limit,
		mode=ImageListCursorMode.RECENTLY,
	),
)

FIRST_LOVE_SPEC: ImageListSpec[tuple[ImageRecord, datetime], DatetimeImageListCursor] = ImageListSpec(
	fetch=lambda repo, cursor, limit: repo.select_first_love(cursor=cursor, limit=limit),
	slice=lambda rows, limit: paginator.slice_with_datetime_tuple_cursor(
		rows,
		limit,
		mode=ImageListCursorMode.FIRST_LOVE,
	),
)

HALL_OF_FAME_SPEC: ImageListSpec[tuple[ImageRecord, datetime], DatetimeImageListCursor] = ImageListSpec(
	fetch=lambda repo, cursor, limit: repo.select_hall_of_fame(cursor=cursor, limit=limit),
	slice=lambda rows, limit: paginator.slice_with_datetime_tuple_cursor(
		rows,
		limit,
		mode=ImageListCursorMode.HALL_OF_FAME,
	),
)

ENGAGED_SPEC: ImageListSpec[tuple[ImageRecord, int], UInt8ImageListCursor] = ImageListSpec(
	fetch=lambda repo, cursor, limit: repo.select_engaged(cursor=cursor, limit=limit),
	slice=lambda rows, limit: paginator.slice_with_uint8_tuple_cursor(
		rows,
		limit,
		mode=ImageListCursorMode.ENGAGED,
	),
)
