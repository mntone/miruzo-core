from collections.abc import Sequence
from datetime import datetime, timedelta, timezone
from typing import Annotated, Any, final

from annotated_types import Len, MaxLen, MinLen
from pydantic import (
	AfterValidator,
	BaseModel,
	BeforeValidator,
	ConfigDict,
	Field,
	PlainSerializer,
	TypeAdapter,
	field_serializer,
	field_validator,
	model_validator,
)

from app.models.enums import ExecutionStatus, ProcessStatus, VisibilityStatus
from app.models.types import IngestIdType, OptionalStrictStr, UtcDateTime

MAX_EXECUTIONS = 5


def _parse_timedelta(value: Any) -> Any:
	if isinstance(value, (int, float)):
		return timedelta(seconds=value)
	return value


def _validate_option_positive_timedelta(value: timedelta | None) -> timedelta | None:
	if value is None:
		return None
	if value.total_seconds() < 0:
		raise ValueError('duration must be greater than or equal to 0')
	return value


def _serialize_optional_timedelta(value: timedelta | None) -> float | None:
	if value is None:
		return None
	else:
		return value.total_seconds()


OptionalPositiveTimeDelta = Annotated[
	timedelta | None,
	BeforeValidator(_parse_timedelta),
	AfterValidator(_validate_option_positive_timedelta),
	PlainSerializer(_serialize_optional_timedelta, return_type=float | None),
]


@final
class Execution(BaseModel):
	model_config = ConfigDict(validate_assignment=True)

	status: ExecutionStatus
	error_type: OptionalStrictStr
	error_message: OptionalStrictStr
	executed_at: datetime
	inspect: OptionalPositiveTimeDelta
	collect: OptionalPositiveTimeDelta
	plan: OptionalPositiveTimeDelta
	execute: OptionalPositiveTimeDelta
	store: OptionalPositiveTimeDelta
	overall: OptionalPositiveTimeDelta

	@field_validator('executed_at')
	@staticmethod
	def validate_executed_at(value: datetime) -> datetime:
		if value.tzinfo is None:
			raise ValueError('timezone-aware datetime required')
		return value

	@field_serializer('executed_at')
	@staticmethod
	def serialize_executed_at(value: datetime) -> float:
		return value.astimezone(timezone.utc).timestamp()


executions_adapter = TypeAdapter(list[Execution])


@final
class Ingest(BaseModel):
	model_config = ConfigDict(validate_assignment=True)

	id: IngestIdType
	process: ProcessStatus = ProcessStatus.PROCESSING
	visibility: VisibilityStatus = VisibilityStatus.PRIVATE
	relative_path: Annotated[str, MinLen(4)]
	fingerprint: Annotated[str, Len(64, 64)]
	ingested_at: UtcDateTime
	captured_at: UtcDateTime
	updated_at: UtcDateTime
	executions: Annotated[Sequence[Execution], Field(default_factory=list), MaxLen(MAX_EXECUTIONS)]

	@field_validator('relative_path')
	@staticmethod
	def validate_relative_path(value: str) -> str:
		if value.startswith('/'):
			raise ValueError('relative_path must not start with "/"')
		return value

	@model_validator(mode='after')
	def validate_datetime(self) -> 'Ingest':
		if self.captured_at > self.ingested_at:
			raise ValueError('captured_at must be less than or equal to ingested_at')
		if self.updated_at < self.ingested_at:
			raise ValueError('updated_at must be greater than or equal to ingested_at')
		return self
