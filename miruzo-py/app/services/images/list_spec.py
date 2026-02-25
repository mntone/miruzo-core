from collections.abc import Sequence
from dataclasses import dataclass
from datetime import datetime
from typing import Callable, Generic, TypeVar, final

import app.services.images.paginator as paginator
from app.models.records import ImageRecord
from app.persist.images.list.protocol import ImageListRepository
from app.services.images.list_types import TCursor

TRow = TypeVar('TRow')


@dataclass(frozen=True, slots=True)
@final
class ImageListSpec(Generic[TRow, TCursor]):
	fetch: Callable[[ImageListRepository, TCursor | None, int], Sequence[TRow]]
	slice: Callable[[Sequence[TRow], int], tuple[Sequence[ImageRecord], TCursor | None]]


LATEST_SPEC: ImageListSpec[ImageRecord, datetime] = ImageListSpec(
	fetch=lambda repo, cursor, limit: repo.select_latest(cursor=cursor, limit=limit),
	slice=paginator.slice_with_latest_cursor,
)

CHRONOLOGICAL_SPEC: ImageListSpec[tuple[ImageRecord, datetime], datetime] = ImageListSpec(
	fetch=lambda repo, cursor, limit: repo.select_chronological(cursor=cursor, limit=limit),
	slice=paginator.slice_with_tuple_cursor,
)

RECENTLY_SPEC: ImageListSpec[tuple[ImageRecord, datetime], datetime] = ImageListSpec(
	fetch=lambda repo, cursor, limit: repo.select_recently(cursor=cursor, limit=limit),
	slice=paginator.slice_with_tuple_cursor,
)

FIRST_LOVE_SPEC: ImageListSpec[tuple[ImageRecord, datetime], datetime] = ImageListSpec(
	fetch=lambda repo, cursor, limit: repo.select_first_love(cursor=cursor, limit=limit),
	slice=paginator.slice_with_tuple_cursor,
)

HALL_OF_FAME_SPEC: ImageListSpec[tuple[ImageRecord, datetime], datetime] = ImageListSpec(
	fetch=lambda repo, cursor, limit: repo.select_hall_of_fame(cursor=cursor, limit=limit),
	slice=paginator.slice_with_tuple_cursor,
)

ENGAGED_SPEC: ImageListSpec[tuple[ImageRecord, int], int] = ImageListSpec(
	fetch=lambda repo, cursor, limit: repo.select_engaged(cursor=cursor, limit=limit),
	slice=paginator.slice_with_tuple_cursor,
)
