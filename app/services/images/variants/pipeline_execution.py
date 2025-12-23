from collections.abc import Iterator
from contextlib import contextmanager
from datetime import datetime, timedelta, timezone
from pathlib import Path
from time import monotonic
from types import TracebackType
from typing import Literal

from PIL import UnidentifiedImageError as PILUnidentifiedImageError
from PIL.Image import DecompressionBombError as PILDecompressionBombError
from sqlalchemy.exc import DataError, IntegrityError, OperationalError

from app.models.enums import ExecutionStatus
from app.models.types import ExecutionEntry
from app.services.images.variants.executors.executor import VariantExecutor
from app.services.images.variants.types import (
	OriginalFile,
	VariantCommitResult,
	VariantPlan,
	VariantPolicy,
)

_ExecutionPhase = Literal['inspect', 'collect', 'plan', 'preprocess', 'commit', 'postprocess', 'store']


class VariantPipelineExecutionSession:
	def __init__(self, executor: VariantExecutor) -> None:
		self._executor = executor

		self._status = ExecutionStatus.SUCCESS
		self._executed_at: datetime | None = None
		self._start_mark: float | None = None
		self._last_mark: float | None = None

		self._error_type: str | None = None
		self._error_message: str | None = None

		self._inspect: timedelta | None = None
		self._collect: timedelta | None = None
		self._plan: timedelta | None = None
		self._preprocess: timedelta | None = None
		# self._commits: list[CommitEntry] = []
		self._commit: timedelta | None = None
		self._postprocess: timedelta | None = None
		self._store: timedelta | None = None
		self._overall: timedelta | None = None

	def __enter__(self) -> 'VariantPipelineExecutionSession':
		self._executed_at = datetime.now(timezone.utc)
		mark = monotonic()
		self._start_mark = mark
		self._last_mark = mark
		return self

	def __exit__(
		self,
		exc_type: type[BaseException] | None,
		exc: BaseException | None,
		tb: TracebackType | None,
	) -> bool | None:
		assert self._start_mark is not None
		self._overall = timedelta(seconds=monotonic() - self._start_mark)

		if exc_type is None:
			self._status = ExecutionStatus.SUCCESS
			return False

		if issubclass(exc_type, PILUnidentifiedImageError):
			self._status = ExecutionStatus.IMAGE_ERROR
		elif issubclass(exc_type, PILDecompressionBombError):
			self._status = ExecutionStatus.IO_ERROR
		elif issubclass(exc_type, (DataError, IntegrityError, OperationalError)):
			self._status = ExecutionStatus.DB_ERROR
		else:
			self._status = ExecutionStatus.UNKNOWN_ERROR

		self._error_type = exc_type.__name__
		self._error_message = exc.__str__()
		return self._status != ExecutionStatus.UNKNOWN_ERROR

	def _mark_phase(self) -> timedelta:
		assert self._last_mark is not None

		current = monotonic()
		duration = timedelta(seconds=current - self._last_mark)
		self._last_mark = current
		return duration

	@contextmanager
	def phase(self, name: _ExecutionPhase) -> Iterator[None]:
		yield None

		duration = self._mark_phase()
		match name:
			case 'inspect':
				self._inspect = duration
			case 'collect':
				self._collect = duration
			case 'plan':
				self._plan = duration
			case 'preprocess':
				self._preprocess = duration
			case 'commit':
				self._commit = duration
			case 'postprocess':
				self._postprocess = duration
			case 'store':
				self._store = duration
			case _:
				raise ValueError(f'Unknown phase: {name}')

	# @contextmanager
	# def commit(self, slot: str) -> Iterator[None]:
	# 	yield None

	# 	duration = self._mark_phase()
	# 	entry = CommitEntry(slot=slot, duration=duration)
	# 	self._commits.append(entry)

	def execute(
		self,
		*,
		media_root: Path,
		file: OriginalFile,
		plan: VariantPlan,
		policy: VariantPolicy,
	) -> Iterator[VariantCommitResult]:
		with self.phase('preprocess'):
			image = self._executor.preprocess(file)

		with self.phase('commit'):
			results = self._executor.commit(
				image,
				media_root=media_root,
				plan=plan,
				policy=policy,
			)

		with self.phase('postprocess'):
			self._executor.postprocess(image)

		return results

	def to_entry(self) -> ExecutionEntry:
		if self._executed_at is None:
			raise RuntimeError('VariantExecutionSession must be used as a context manager')

		entry = ExecutionEntry(
			status=self._status,
			error_type=self._error_type,
			error_message=self._error_message,
			executed_at=self._executed_at,
			inspect=self._inspect,
			collect=self._collect,
			plan=self._plan,
			preprocess=self._preprocess,
			commit=self._commit,
			# commits=self._commits if len(self._commits) > 0 else None,
			postprocess=self._postprocess,
			store=self._store,
			overall=self._overall,
		)
		return entry
