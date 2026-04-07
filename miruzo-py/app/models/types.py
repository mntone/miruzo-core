from datetime import datetime, timedelta, timezone
from typing import Annotated, TypedDict, final

from annotated_types import Interval, Len
from pydantic import AfterValidator, BeforeValidator, Field, StrictStr

from app.config import constants as c
from app.config.constants import MAX_IMAGE_HEIGHT, MAX_IMAGE_WIDTH

OptionalStrictStr = StrictStr | None

IngestIdType = Annotated[int, Interval(ge=c.INGEST_ID_MINIMUM, le=c.INGEST_ID_MAXIMUM)]


def _validate_relative_path(value: str) -> str:
	if value.startswith('/') or value.startswith('..'):
		raise ValueError('relative_path must not start with "/" or ".."')
	return value


RelativePathType = Annotated[str, Len(5, 255), BeforeValidator(_validate_relative_path)]


def _to_utc(value: datetime) -> datetime:
	if value.tzinfo is None:
		value = value.replace(tzinfo=timezone.utc)
	else:
		value = value.astimezone(timezone.utc)
	return value


UtcDateTime = Annotated[datetime, AfterValidator(_to_utc)]


@final
class CommitEntry(TypedDict):
	slot: str
	duration: Annotated[timedelta, Field(ge=0)]


@final
class VariantEntry(TypedDict):
	rel: str
	layer_id: Annotated[int, Field(ge=0, le=9)]
	format: Annotated[str, Field(min_length=3, max_length=8)]
	codecs: Annotated[str | None, Field(default=None)]
	bytes: Annotated[int, Field(ge=1)]
	width: Annotated[int, Field(ge=1, le=MAX_IMAGE_WIDTH)]
	height: Annotated[int, Field(ge=1, le=MAX_IMAGE_HEIGHT)]
	quality: Annotated[int | None, Field(default=None, ge=1, le=100)]
