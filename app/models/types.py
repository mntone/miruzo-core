from collections.abc import Sequence
from datetime import datetime, timedelta, timezone
from typing import Annotated, Any, TypedDict, final

from pydantic import Field
from sqlalchemy.types import JSON, TypeDecorator

from app.models.enums import ExecutionStatus

# @final
# class _CommitPlainJson(TypedDict):
# 	slot: str
# 	duration: float


@final
class _ExecutionPlainJson(TypedDict):
	status: ExecutionStatus
	error_type: str | None
	error_message: str | None
	executed_at: float
	"""unix epoch seconds (UTC)"""
	inspect: float | None
	collect: float | None
	plan: float | None
	execute: float | None
	# commits: Sequence[_CommitPlainJson] | None
	store: float | None
	overall: float | None


@final
class CommitEntry(TypedDict):
	slot: str
	duration: Annotated[timedelta, Field(ge=0)]


@final
class ExecutionEntry(TypedDict):
	status: ExecutionStatus
	error_type: str | None
	error_message: str | None
	executed_at: datetime
	inspect: Annotated[timedelta | None, Field(ge=0)]
	collect: Annotated[timedelta | None, Field(ge=0)]
	plan: Annotated[timedelta | None, Field(ge=0)]
	execute: Annotated[timedelta | None, Field(ge=0)]
	# commits: Annotated[Sequence[CommitEntry] | None, Field(min_length=1)]
	store: Annotated[timedelta | None, Field(ge=0)]
	overall: Annotated[timedelta | None, Field(ge=0)]


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
				error_type=v['error_type'],
				error_message=v['error_message'],
				executed_at=(
					v['executed_at'].replace(tzinfo=timezone.utc)
					if v['executed_at'].tzinfo is None
					else v['executed_at'].astimezone(timezone.utc)
				).timestamp(),
				inspect=ExecutionsJSON.to_seconds(v['inspect']),
				collect=ExecutionsJSON.to_seconds(v['collect']),
				plan=ExecutionsJSON.to_seconds(v['plan']),
				execute=ExecutionsJSON.to_seconds(v.get('execute')),
				# commits=(
				# 	[
				# 		_CommitPlainJson(slot=c['slot'], duration=c['duration'].total_seconds())
				# 		for c in v['commits']
				# 	]
				# 	if v['commits']
				# 	else None
				# ),
				store=ExecutionsJSON.to_seconds(v['store']),
				overall=ExecutionsJSON.to_seconds(v.get('overall')),
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
				error_type=v['error_type'],
				error_message=v['error_message'],
				executed_at=datetime.fromtimestamp(v['executed_at'], tz=timezone.utc),
				inspect=ExecutionsJSON.to_timedelta(v['inspect']),
				collect=ExecutionsJSON.to_timedelta(v['collect']),
				plan=ExecutionsJSON.to_timedelta(v['plan']),
				execute=ExecutionsJSON.to_timedelta(v.get('execute')),
				# commits=(
				# 	[
				# 		CommitEntry(slot=c['slot'], duration=timedelta(seconds=c['duration']))
				# 		for c in v['commits']
				# 	]
				# 	if v['commits']
				# 	else None
				# ),
				store=ExecutionsJSON.to_timedelta(v['store']),
				overall=ExecutionsJSON.to_timedelta(v.get('overall')),
			)
			for v in value
		]

	@staticmethod
	def to_seconds(delta: timedelta | None) -> float | None:
		if delta is None:
			return None
		else:
			return delta.total_seconds()

	@staticmethod
	def to_timedelta(seconds: float | None) -> timedelta | None:
		if seconds is None:
			return None
		else:
			return timedelta(seconds=seconds)


@final
class VariantEntry(TypedDict):
	rel: str
	layer_id: Annotated[int, Field(ge=0, le=9)]
	format: Annotated[str, Field(ge=3, le=8)]
	codecs: Annotated[str | None, Field(default=None)]
	bytes: Annotated[int, Field(ge=1)]
	width: Annotated[int, Field(ge=1, le=10240)]
	height: Annotated[int, Field(ge=1, le=10240)]
	quality: Annotated[int | None, Field(default=None, ge=1, le=100)]
