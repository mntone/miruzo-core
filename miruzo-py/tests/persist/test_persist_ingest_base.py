import json
from datetime import datetime, timedelta, timezone
from typing import Any, Generator

import pytest
from sqlalchemy import create_engine
from sqlalchemy.exc import NoResultFound
from sqlalchemy.orm import Session

from tests.persist.utils import get_ingest_dto, get_ingest_row

from app.config.constants import EXECUTION_MAXIMUM
from app.databases.metadata import metadata
from app.models.enums import ExecutionStatus
from app.models.ingest import Execution
from app.persist.ingests.base import _IngestRepositoryBaseImpl
from app.persist.ingests.protocol import IngestAppendExecutionInput, IngestCreateInput


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


def test_append_execution_replaces_previous_success(session: Session) -> None:
	now = datetime.now(timezone.utc)
	repo = _IngestRepositoryBaseImpl(session)
	ingest_id = repo.create(
		IngestCreateInput(
			relative_path='l0orig/foo.webp',
			fingerprint='0' * 64,
			ingested_at=now,
			captured_at=now,
		),
	)

	repo.append_execution(
		IngestAppendExecutionInput(
			ingest_id=ingest_id,
			updated_at=now,
			execution=_build_execution(ExecutionStatus.UNKNOWN_ERROR),
		),
	)
	row = get_ingest_row(session, ingest_id=ingest_id)
	assert row is not None
	assert row['executions'] is not None
	assert [e['status'] for e in row['executions']] == [ExecutionStatus.UNKNOWN_ERROR]

	repo.append_execution(
		IngestAppendExecutionInput(
			ingest_id=ingest_id,
			updated_at=now,
			execution=_build_execution(ExecutionStatus.SUCCESS),
		),
	)
	row = get_ingest_row(session, ingest_id=ingest_id)
	assert row is not None
	assert row['executions'] is not None
	assert [e['status'] for e in row['executions']] == [
		ExecutionStatus.UNKNOWN_ERROR,
		ExecutionStatus.SUCCESS,
	]

	repo.append_execution(
		IngestAppendExecutionInput(
			ingest_id=ingest_id,
			updated_at=now,
			execution=_build_execution(ExecutionStatus.SUCCESS),
		),
	)
	row = get_ingest_row(session, ingest_id=ingest_id)
	assert row is not None
	assert row['executions'] is not None
	assert [e['status'] for e in row['executions']] == [
		ExecutionStatus.UNKNOWN_ERROR,
		ExecutionStatus.SUCCESS,
	]


def test_append_execution_trims_to_maximum(session: Session) -> None:
	now = datetime.now(timezone.utc)
	repo = _IngestRepositoryBaseImpl(session)
	ingest_id = repo.create(
		IngestCreateInput(
			relative_path='l0orig/foo.webp',
			fingerprint='1' * 64,
			ingested_at=now,
			captured_at=now,
		),
	)
	assert ingest_id == 1

	row = get_ingest_row(session, ingest_id=ingest_id)
	assert row is not None
	assert len(row.executions) == 0

	for i in range(EXECUTION_MAXIMUM + 3):
		repo.append_execution(
			IngestAppendExecutionInput(
				ingest_id=ingest_id,
				updated_at=now,
				execution=_build_execution(ExecutionStatus.UNKNOWN_ERROR, offset=i),
			),
		)

		row = get_ingest_row(session, ingest_id=ingest_id)
		assert len(row.executions) == min(1 + i, EXECUTION_MAXIMUM)

	row = get_ingest_dto(session, ingest_id=ingest_id)
	assert row is not None
	assert len(row.executions) == EXECUTION_MAXIMUM

	last_execution = row.executions[-1]
	assert last_execution.inspect is not None
	assert last_execution.inspect.total_seconds() == (EXECUTION_MAXIMUM + 2) + 10
	assert last_execution.collect is not None
	assert last_execution.collect.total_seconds() == (EXECUTION_MAXIMUM + 2) + 20
	assert last_execution.plan is not None
	assert last_execution.plan.total_seconds() == (EXECUTION_MAXIMUM + 2) + 30
	assert last_execution.execute is not None
	assert last_execution.execute.total_seconds() == (EXECUTION_MAXIMUM + 2) + 50
	assert last_execution.store is not None
	assert last_execution.store.total_seconds() == (EXECUTION_MAXIMUM + 2) + 100
	assert last_execution.overall is not None
	assert last_execution.overall.total_seconds() == (EXECUTION_MAXIMUM + 2) + 120


def test_append_execution_returns_none_for_missing_ingest(session: Session) -> None:
	now = datetime.now(timezone.utc)
	repo = _IngestRepositoryBaseImpl(session)
	with pytest.raises(NoResultFound, match='No row was found when one was required'):
		repo.append_execution(
			IngestAppendExecutionInput(
				ingest_id=999,
				updated_at=now,
				execution=_build_execution(ExecutionStatus.SUCCESS),
			),
		)
