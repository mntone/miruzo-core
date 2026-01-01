from datetime import datetime, timedelta, timezone
from typing import Any, Generator

import pytest
from sqlmodel import Session, SQLModel, create_engine

from app.config.constants import EXECUTION_MAXIMUM
from app.models.enums import ExecutionStatus
from app.models.types import ExecutionEntry
from app.services.ingests.repository.base import IngestRepository


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	SQLModel.metadata.create_all(engine)
	with Session(engine) as session:
		yield session


def _build_execution(status: ExecutionStatus, *, offset: int = 0) -> ExecutionEntry:
	now = datetime(2026, 1, 1, tzinfo=timezone.utc)
	return {
		'status': status,
		'error_type': None,
		'error_message': None,
		'executed_at': now,
		'inspect': timedelta(seconds=offset + 10),
		'collect': timedelta(seconds=offset + 20),
		'plan': timedelta(seconds=offset + 30),
		'execute': timedelta(seconds=offset + 50),
		#'commits': [CommitEntry(slot='l1w320', duration=timedelta(seconds=offset + 50))],
		'store': timedelta(seconds=offset + 100),
		'overall': timedelta(seconds=offset + 120),
	}


def test_append_execution_replaces_previous_success(session: Session) -> None:
	repo = IngestRepository(session)
	now = datetime.now(timezone.utc)
	ingest = repo.create_ingest(
		relative_path='l0orig/foo.webp',
		fingerprint='0' * 64,
		ingested_at=now,
		captured_at=now,
	)

	ingest = repo.append_execution(ingest.id, _build_execution(ExecutionStatus.UNKNOWN_ERROR))
	assert ingest is not None
	assert ingest.executions is not None
	assert [e['status'] for e in ingest.executions] == [ExecutionStatus.UNKNOWN_ERROR]

	ingest = repo.append_execution(ingest.id, _build_execution(ExecutionStatus.SUCCESS))
	assert ingest is not None
	assert ingest.executions is not None
	assert [e['status'] for e in ingest.executions] == [
		ExecutionStatus.UNKNOWN_ERROR,
		ExecutionStatus.SUCCESS,
	]

	ingest = repo.append_execution(ingest.id, _build_execution(ExecutionStatus.SUCCESS))
	assert ingest is not None
	assert ingest.executions is not None
	assert [e['status'] for e in ingest.executions] == [
		ExecutionStatus.UNKNOWN_ERROR,
		ExecutionStatus.SUCCESS,
	]


def test_append_execution_trims_to_maximum(session: Session) -> None:
	repo = IngestRepository(session)
	now = datetime.now(timezone.utc)
	ingest = repo.create_ingest(
		relative_path='l0orig/foo.webp',
		fingerprint='1' * 64,
		ingested_at=now,
		captured_at=now,
	)
	assert ingest is not None

	for i in range(EXECUTION_MAXIMUM + 3):
		ingest = repo.append_execution(
			ingest.id,
			_build_execution(ExecutionStatus.UNKNOWN_ERROR, offset=i),
		)
		assert ingest is not None
		assert ingest.executions is not None
		assert len(ingest.executions) == min(1 + i, EXECUTION_MAXIMUM)

	assert ingest is not None
	assert ingest.executions is not None
	assert len(ingest.executions) == EXECUTION_MAXIMUM
	assert ingest.executions[-1]['inspect'] is not None
	assert ingest.executions[-1]['inspect'].total_seconds() == (EXECUTION_MAXIMUM + 2) + 10
	assert ingest.executions[-1]['collect'] is not None
	assert ingest.executions[-1]['collect'].total_seconds() == (EXECUTION_MAXIMUM + 2) + 20
	assert ingest.executions[-1]['plan'] is not None
	assert ingest.executions[-1]['plan'].total_seconds() == (EXECUTION_MAXIMUM + 2) + 30
	assert ingest.executions[-1]['execute'] is not None
	assert ingest.executions[-1]['execute'].total_seconds() == (EXECUTION_MAXIMUM + 2) + 50
	assert ingest.executions[-1]['store'] is not None
	assert ingest.executions[-1]['store'].total_seconds() == (EXECUTION_MAXIMUM + 2) + 100
	assert ingest.executions[-1]['overall'] is not None
	assert ingest.executions[-1]['overall'].total_seconds() == (EXECUTION_MAXIMUM + 2) + 120


def test_append_execution_returns_none_for_missing_ingest(session: Session) -> None:
	repo = IngestRepository(session)

	result = repo.append_execution(999, _build_execution(ExecutionStatus.SUCCESS))

	assert result is None
