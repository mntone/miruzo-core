from collections.abc import Generator
from datetime import datetime, timezone
from typing import Any

import pytest
from sqlmodel import Session, SQLModel, create_engine

from app.persist.jobs.sqlite import SQLiteJobRepository


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	SQLModel.metadata.create_all(engine)
	with Session(engine) as session:
		yield session


def test_get_or_create_returns_existing(session: Session) -> None:
	repo = SQLiteJobRepository(session)

	job = repo.get_or_create('score_decay')

	again = repo.get_or_create('score_decay')

	assert again is job


def test_mark_started_and_finished(session: Session) -> None:
	repo = SQLiteJobRepository(session)
	job = repo.get_or_create('score_decay')

	started_at = datetime(2026, 1, 2, 6, 0, tzinfo=timezone.utc)
	finished_at = datetime(2026, 1, 2, 6, 1, tzinfo=timezone.utc)

	repo.mark_started(job, started_at=started_at)
	repo.mark_finished(job.name, finished_at=finished_at)

	assert job.started_at == started_at
	assert job.finished_at == finished_at
