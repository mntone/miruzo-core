from dataclasses import dataclass
from datetime import datetime
from enum import Enum
from typing import TypeAlias, TypeVar


class ImageListCursorMode(int, Enum):
	LATEST = 1
	CHRONOLOGICAL = 2
	RECENTLY = 3
	FIRST_LOVE = 4
	HALL_OF_FAME = 5
	ENGAGED = 6


@dataclass(frozen=True, slots=True)
class DatetimeImageListCursor:
	mode: ImageListCursorMode
	value: datetime
	ingest_id: int


@dataclass(frozen=True, slots=True)
class UInt8ImageListCursor:
	mode: ImageListCursorMode
	value: int
	ingest_id: int


ImageListCursor: TypeAlias = DatetimeImageListCursor | UInt8ImageListCursor

TImageListCursor = TypeVar('TImageListCursor', DatetimeImageListCursor, UInt8ImageListCursor)
