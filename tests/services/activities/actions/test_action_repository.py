from collections.abc import Generator
from datetime import datetime, timedelta, timezone
from typing import Any

import pytest
from sqlalchemy.exc import MultipleResultsFound
from sqlmodel import Session, SQLModel, create_engine

from tests.services.images.utils import add_ingest_record

from app.models.enums import ActionKind
from app.services.activities.actions.repository import ActionRepository


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	SQLModel.metadata.create_all(engine)
	with Session(engine) as session:
		yield session


def test_select_by_ingest_id_orders_by_occurred_at(session: Session) -> None:
	ingest = add_ingest_record(session, 1)
	repo = ActionRepository(session)

	repo.insert(
		ingest.id,
		kind=ActionKind.VIEW,
		occurred_at=datetime(2024, 1, 2, tzinfo=timezone.utc),
	)
	repo.insert(
		ingest.id,
		kind=ActionKind.MEMO,
		occurred_at=datetime(2024, 1, 1, tzinfo=timezone.utc),
	)

	items = repo.select_by_ingest_id(ingest.id)

	assert [item.kind for item in items] == [
		ActionKind.MEMO,
		ActionKind.VIEW,
	]


def test_select_by_ingest_id_returns_empty_when_missing(session: Session) -> None:
	ingest = add_ingest_record(session, 1)
	repo = ActionRepository(session)

	items = repo.select_by_ingest_id(ingest.id)

	assert items == []


def test_select_latest_one_respects_time_window(session: Session) -> None:
	ingest = add_ingest_record(session, 1)
	repo = ActionRepository(session)

	since = datetime(2026, 1, 1, tzinfo=timezone.utc)
	until = datetime(2026, 1, 2, tzinfo=timezone.utc)

	repo.insert(
		ingest.id,
		kind=ActionKind.DECAY,
		occurred_at=since,
	)
	repo.insert(
		ingest.id,
		kind=ActionKind.DECAY,
		occurred_at=until,
	)
	repo.insert(
		ingest.id,
		kind=ActionKind.DECAY,
		occurred_at=datetime(2026, 1, 1, 12, 0, tzinfo=timezone.utc),
	)

	matched = repo.select_latest_one(
		ingest.id,
		kind=ActionKind.DECAY,
		since_occurred_at=since,
		until_occurred_at=until,
	)
	assert matched is not None
	assert matched.occurred_at < until

	boundary_match = repo.select_latest_one(
		ingest.id,
		kind=ActionKind.DECAY,
		since_occurred_at=until,
		until_occurred_at=until + timedelta(days=1),
	)
	assert boundary_match is not None
	assert boundary_match.occurred_at == until


def test_select_latest_one_returns_none_when_missing(session: Session) -> None:
	ingest = add_ingest_record(session, 1)
	repo = ActionRepository(session)

	since = datetime(2024, 1, 1, tzinfo=timezone.utc)
	until = datetime(2024, 1, 2, tzinfo=timezone.utc)

	result = repo.select_latest_one(
		ingest.id,
		kind=ActionKind.DECAY,
		since_occurred_at=since,
		until_occurred_at=until,
	)
	assert result is None


def test_select_latest_one_raises_when_require_unique(session: Session) -> None:
	ingest = add_ingest_record(session, 1)
	repo = ActionRepository(session)

	at_time = datetime(2026, 1, 1, tzinfo=timezone.utc)
	repo.insert(
		ingest.id,
		kind=ActionKind.DECAY,
		occurred_at=at_time,
	)
	repo.insert(
		ingest.id,
		kind=ActionKind.DECAY,
		occurred_at=at_time + timedelta(hours=1),
	)

	with pytest.raises(MultipleResultsFound):
		repo.select_latest_one(
			ingest.id,
			kind=ActionKind.DECAY,
			since_occurred_at=at_time,
			require_unique=True,
		)


def test_select_latest_one_prefers_latest_id_on_tie(session: Session) -> None:
	ingest = add_ingest_record(session, 1)
	repo = ActionRepository(session)

	at_time = datetime(2026, 1, 1, tzinfo=timezone.utc)
	first = repo.insert(
		ingest.id,
		kind=ActionKind.DECAY,
		occurred_at=at_time,
	)
	second = repo.insert(
		ingest.id,
		kind=ActionKind.DECAY,
		occurred_at=at_time,
	)

	result = repo.select_latest_one(
		ingest.id,
		kind=ActionKind.DECAY,
		since_occurred_at=at_time,
	)
	assert result is not None
	assert first.id is not None
	assert second.id is not None
	assert result.id == max(first.id, second.id)
