from collections.abc import Iterator
from dataclasses import dataclass
from datetime import datetime, timedelta, timezone
from typing import final

import pytest
from sqlalchemy.exc import NoResultFound

from tests.persist.utils import get_ingest_dto, get_ingest_row

from app.config.environments import DatabaseBackend
from app.models.enums import ExecutionStatus
from app.models.ingest import MAX_EXECUTIONS, Execution
from app.persist.ingests.factory import _create_ingest_repository_from_backend
from app.persist.ingests.protocol import IngestAppendExecutionInput, IngestCreateInput, IngestRepository


@pytest.fixture()
def ingest_repo(request: pytest.FixtureRequest) -> Iterator[IngestRepository]:
	with request.getfixturevalue('session') as session:
		yield _create_ingest_repository_from_backend(
			session,
			backend=DatabaseBackend.SQLITE,
			max_executions=MAX_EXECUTIONS,
		)


def _build_execution(status: ExecutionStatus, *, offset: int = 0) -> Execution:
	now = datetime(2026, 1, 1, tzinfo=timezone.utc)
	return Execution(
		status=status,
		error_type=None,
		error_message=None,
		executed_at=now,
		inspect=timedelta(seconds=offset + 10),
		collect=timedelta(seconds=offset + 20),
		plan=timedelta(seconds=offset + 30),
		execute=timedelta(seconds=offset + 50),
		# commits=[CommitEntry(slot='l1w320', duration=timedelta(seconds=offset + 50))],
		store=timedelta(seconds=offset + 100),
		overall=timedelta(seconds=offset + 120),
	)


def _assert_execution(execution: Execution, *, offset: int = 0) -> None:
	assert execution.inspect is not None
	assert execution.inspect.total_seconds() == offset + 10
	assert execution.collect is not None
	assert execution.collect.total_seconds() == offset + 20
	assert execution.plan is not None
	assert execution.plan.total_seconds() == offset + 30
	assert execution.execute is not None
	assert execution.execute.total_seconds() == offset + 50
	assert execution.store is not None
	assert execution.store.total_seconds() == offset + 100
	assert execution.overall is not None
	assert execution.overall.total_seconds() == offset + 120


def test_append_execution_replaces_previous_success(ingest_repo: IngestRepository) -> None:
	now = datetime.now(timezone.utc)
	ingest_id = ingest_repo.create(
		IngestCreateInput(
			relative_path='l0orig/foo.webp',
			fingerprint='0' * 64,
			ingested_at=now,
			captured_at=now,
		),
	)

	ingest_repo.append_execution(
		IngestAppendExecutionInput(
			ingest_id=ingest_id,
			updated_at=now,
			execution=_build_execution(ExecutionStatus.UNKNOWN_ERROR),
		),
	)
	row = get_ingest_row(ingest_repo, ingest_id=ingest_id)
	assert row is not None
	assert row['executions'] is not None
	assert [e['status'] for e in row['executions']] == [ExecutionStatus.UNKNOWN_ERROR]

	ingest_repo.append_execution(
		IngestAppendExecutionInput(
			ingest_id=ingest_id,
			updated_at=now,
			execution=_build_execution(ExecutionStatus.SUCCESS),
		),
	)
	row = get_ingest_row(ingest_repo, ingest_id=ingest_id)
	assert row is not None
	assert row['executions'] is not None
	assert [e['status'] for e in row['executions']] == [
		ExecutionStatus.UNKNOWN_ERROR,
		ExecutionStatus.SUCCESS,
	]

	ingest_repo.append_execution(
		IngestAppendExecutionInput(
			ingest_id=ingest_id,
			updated_at=now,
			execution=_build_execution(ExecutionStatus.SUCCESS),
		),
	)
	row = get_ingest_row(ingest_repo, ingest_id=ingest_id)
	assert row is not None
	assert row['executions'] is not None
	assert [e['status'] for e in row['executions']] == [
		ExecutionStatus.UNKNOWN_ERROR,
		ExecutionStatus.SUCCESS,
	]


@final
@dataclass(slots=True)
class IngestAppendExecutionContext:
	repo: IngestRepository
	id: int
	max_executions: int

	def append_execution_and_assert(
		self,
		updated_at: datetime,
		*,
		offset: int,
		status: ExecutionStatus = ExecutionStatus.UNKNOWN_ERROR,
	) -> None:
		self.repo.append_execution(
			IngestAppendExecutionInput(
				ingest_id=self.id,
				updated_at=updated_at,
				execution=_build_execution(status, offset=offset),
			),
		)

		row = get_ingest_row(self.repo, ingest_id=self.id)
		assert len(row.executions) == min(1 + offset, self.max_executions)


def test_append_execution_trims_to_maximum(ingest_repo: IngestRepository) -> None:
	now = datetime.now(timezone.utc)
	ingest_id = ingest_repo.create(
		IngestCreateInput(
			relative_path='l0orig/foo.webp',
			fingerprint='1' * 64,
			ingested_at=now,
			captured_at=now,
		),
	)
	assert ingest_id == 1

	context = IngestAppendExecutionContext(
		repo=ingest_repo,
		id=ingest_id,
		max_executions=MAX_EXECUTIONS,
	)

	# Initial state: no executions yet.
	row = get_ingest_row(ingest_repo, ingest_id=ingest_id)
	assert row is not None
	assert len(row.executions) == 0

	# Add one success execution first.
	context.append_execution_and_assert(now, offset=0, status=ExecutionStatus.SUCCESS)

	# Add enough error executions to exceed the current cap.
	EXTRA_EXECUTION_COUNT = 3
	for i in range(1, MAX_EXECUTIONS + EXTRA_EXECUTION_COUNT):
		context.append_execution_and_assert(now, offset=i)

	dto1 = get_ingest_dto(ingest_repo, ingest_id=ingest_id)
	assert dto1 is not None
	assert len(dto1.executions) == MAX_EXECUTIONS
	assert dto1.executions[0].status == ExecutionStatus.UNKNOWN_ERROR
	_assert_execution(dto1.executions[0], offset=EXTRA_EXECUTION_COUNT)
	_assert_execution(dto1.executions[-1], offset=MAX_EXECUTIONS + EXTRA_EXECUTION_COUNT - 1)

	# White-box test: lower the internal cap at runtime and verify shrink behavior.
	NEW_MAX_EXECUTIONS = 3
	context.max_executions = NEW_MAX_EXECUTIONS
	ingest_repo._max_executions = NEW_MAX_EXECUTIONS  # pyright: ignore[reportAttributeAccessIssue]

	for i in range(
		MAX_EXECUTIONS + EXTRA_EXECUTION_COUNT,
		MAX_EXECUTIONS + EXTRA_EXECUTION_COUNT + 3,
	):
		context.append_execution_and_assert(now, offset=i)

	dto2 = get_ingest_dto(ingest_repo, ingest_id=ingest_id)
	assert dto2 is not None
	assert len(dto2.executions) == NEW_MAX_EXECUTIONS
	_assert_execution(dto2.executions[0], offset=MAX_EXECUTIONS + EXTRA_EXECUTION_COUNT)
	_assert_execution(dto2.executions[1], offset=MAX_EXECUTIONS + EXTRA_EXECUTION_COUNT + 1)
	_assert_execution(dto2.executions[2], offset=MAX_EXECUTIONS + EXTRA_EXECUTION_COUNT + 2)


def test_append_execution_returns_none_for_missing_ingest(ingest_repo: IngestRepository) -> None:
	now = datetime.now(timezone.utc)
	with pytest.raises(NoResultFound, match='No row was found when one was required'):
		ingest_repo.append_execution(
			IngestAppendExecutionInput(
				ingest_id=999,
				updated_at=now,
				execution=_build_execution(ExecutionStatus.SUCCESS),
			),
		)
