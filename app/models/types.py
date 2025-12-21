from collections.abc import Sequence
from datetime import datetime, timedelta, timezone
from typing import Annotated, Any, TypedDict, final

from pydantic import Field
from sqlalchemy.types import JSON, TypeDecorator

from app.models.enums import ExecutionStatus


@final
class _CommitPlainJson(TypedDict):
	slot: str
	duration: float


@final
class _ExecutionPlainJson(TypedDict):
	status: ExecutionStatus
	executed_at: float
	"""unix epoch seconds (UTC)"""
	collect: float
	plan: float
	preprocess: float
	commit: Sequence[_CommitPlainJson] | None
	store: float | None


@final
class CommitEntry(TypedDict):
	slot: str
	duration: Annotated[timedelta, Field(ge=0)]


@final
class ExecutionEntry(TypedDict):
	status: ExecutionStatus
	executed_at: datetime
	collect: Annotated[timedelta, Field(ge=0)]
	plan: Annotated[timedelta, Field(ge=0)]
	preprocess: Annotated[timedelta, Field(ge=0)]
	commit: Annotated[Sequence[CommitEntry] | None, Field(min_length=1)]
	store: Annotated[timedelta | None, Field(ge=0)]


class ExecutionsJSON(TypeDecorator[Sequence[ExecutionEntry] | None]):
	impl = JSON
	cache_ok = True

	def process_bind_param(
		self,
		value: Sequence[ExecutionEntry] | None,
		dialect: Any,  # noqa: ARG002
	) -> Sequence[_ExecutionPlainJson] | None:
		if value is None:
			return None
		return [
			_ExecutionPlainJson(
				status=v['status'],
				executed_at=(
					v['executed_at'].replace(tzinfo=timezone.utc)
					if v['executed_at'].tzinfo is None
					else v['executed_at'].astimezone(timezone.utc)
				).timestamp(),
				collect=v['collect'].total_seconds(),
				plan=v['plan'].total_seconds(),
				preprocess=v['preprocess'].total_seconds(),
				commit=(
					[
						_CommitPlainJson(slot=c['slot'], duration=c['duration'].total_seconds())
						for c in v['commit']
					]
					if v['commit']
					else None
				),
				store=v['store'].total_seconds() if v['store'] is not None else None,
			)
			for v in value
		]

	def process_result_value(
		self,
		value: Sequence[_ExecutionPlainJson] | None,
		dialect: Any,  # noqa: ARG002
	) -> Sequence[ExecutionEntry] | None:
		if value is None:
			return None
		return [
			ExecutionEntry(
				status=ExecutionStatus(v['status']),
				executed_at=datetime.fromtimestamp(v['executed_at'], tz=timezone.utc),
				collect=timedelta(seconds=v['collect']),
				plan=timedelta(seconds=v['plan']),
				preprocess=timedelta(seconds=v['preprocess']),
				commit=(
					[
						CommitEntry(slot=c['slot'], duration=timedelta(seconds=c['duration']))
						for c in v['commit']
					]
					if v['commit']
					else None
				),
				store=timedelta(seconds=v['store']) if v['store'] is not None else None,
			)
			for v in value
		]
