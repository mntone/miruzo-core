from datetime import datetime, timedelta, timezone
from typing import Any, Generator

import pytest
from sqlmodel import Session, SQLModel, create_engine

from tests.services.activities.stats.factory import add_stats_record
from tests.services.images.utils import add_image_record

from app.domain.images.cursor import (
	DatetimeImageListCursor,
	ImageListCursorMode,
	UInt8ImageListCursor,
)
from app.persist.images.list.base import BaseImageListRepository


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	SQLModel.metadata.create_all(engine)
	with Session(engine) as session:
		yield session


def test_select_latest_orders_desc_and_sets_cursor(session: Session) -> None:
	now = datetime.now(timezone.utc)
	first = add_image_record(session, 4, ingested_at=now)
	second = add_image_record(session, 3, ingested_at=now - timedelta(hours=1))
	second_tie = add_image_record(session, 2, ingested_at=now - timedelta(hours=1))
	third = add_image_record(session, 1, ingested_at=now - timedelta(hours=2))

	repository = BaseImageListRepository(session, engaged_score_threshold=160)

	rows = repository.select_latest(cursor=None, limit=2)
	assert [row.ingest_id for row in rows] == [first.ingest_id, second.ingest_id]

	cursor = DatetimeImageListCursor(
		mode=ImageListCursorMode.LATEST,
		value=rows[-1].ingested_at,
		ingest_id=rows[-1].ingest_id,
	)
	next_rows = repository.select_latest(cursor=cursor, limit=2)
	assert [row.ingest_id for row in next_rows] == [second_tie.ingest_id, third.ingest_id]


def test_select_chronological_orders_by_captured_at(session: Session) -> None:
	now = datetime(2024, 1, 2, tzinfo=timezone.utc)
	first = add_image_record(session, 4, captured_at=now)
	second = add_image_record(session, 3, captured_at=now - timedelta(hours=1))
	second_tie = add_image_record(session, 2, captured_at=now - timedelta(hours=1))
	third = add_image_record(session, 1, captured_at=now - timedelta(hours=2))

	repository = BaseImageListRepository(session, engaged_score_threshold=160)

	rows = repository.select_chronological(cursor=None, limit=2)
	assert [row[0].ingest_id for row in rows] == [first.ingest_id, second.ingest_id]

	cursor = DatetimeImageListCursor(
		mode=ImageListCursorMode.CHRONOLOGICAL,
		value=rows[-1][1],
		ingest_id=rows[-1][0].ingest_id,
	)
	next_rows = repository.select_chronological(cursor=cursor, limit=2)
	assert [row[0].ingest_id for row in next_rows] == [second_tie.ingest_id, third.ingest_id]


def test_select_recently_orders_by_last_viewed_at(session: Session) -> None:
	now = datetime(2024, 1, 2, tzinfo=timezone.utc)
	first = add_image_record(session, 4, ingested_at=now)
	second = add_image_record(session, 3, ingested_at=now - timedelta(hours=1))
	second_tie = add_image_record(session, 2, ingested_at=now - timedelta(hours=2))
	third = add_image_record(session, 1, ingested_at=now - timedelta(hours=3))

	add_stats_record(
		session,
		first.ingest_id,
		view_count=1,
		last_viewed_at=now,
	)
	add_stats_record(
		session,
		second.ingest_id,
		view_count=1,
		last_viewed_at=now - timedelta(days=1),
	)
	add_stats_record(
		session,
		second_tie.ingest_id,
		view_count=1,
		last_viewed_at=now - timedelta(days=1),
	)
	add_stats_record(
		session,
		third.ingest_id,
		view_count=1,
		last_viewed_at=now - timedelta(days=2),
	)

	repository = BaseImageListRepository(session, engaged_score_threshold=160)

	rows = repository.select_recently(cursor=None, limit=2)
	assert [row[0].ingest_id for row in rows] == [first.ingest_id, second.ingest_id]

	cursor = DatetimeImageListCursor(
		mode=ImageListCursorMode.RECENTLY,
		value=rows[-1][1],
		ingest_id=rows[-1][0].ingest_id,
	)
	next_rows = repository.select_recently(cursor=cursor, limit=2)
	assert [row[0].ingest_id for row in next_rows] == [second_tie.ingest_id, third.ingest_id]


def test_select_first_love_orders_by_first_loved_at(session: Session) -> None:
	now = datetime(2024, 1, 2, tzinfo=timezone.utc)
	first = add_image_record(session, 4, ingested_at=now)
	second = add_image_record(session, 3, ingested_at=now - timedelta(hours=1))
	second_tie = add_image_record(session, 2, ingested_at=now - timedelta(hours=2))
	third = add_image_record(session, 1, ingested_at=now - timedelta(hours=3))

	add_stats_record(
		session,
		first.ingest_id,
		view_count=1,
		first_loved_at=now,
	)
	add_stats_record(
		session,
		second.ingest_id,
		view_count=1,
		first_loved_at=now - timedelta(days=1),
	)
	add_stats_record(
		session,
		second_tie.ingest_id,
		view_count=1,
		first_loved_at=now - timedelta(days=1),
	)
	add_stats_record(
		session,
		third.ingest_id,
		view_count=1,
		first_loved_at=now - timedelta(days=2),
	)

	repository = BaseImageListRepository(session, engaged_score_threshold=160)

	rows = repository.select_first_love(cursor=None, limit=2)
	assert [row[0].ingest_id for row in rows] == [first.ingest_id, second.ingest_id]

	cursor = DatetimeImageListCursor(
		mode=ImageListCursorMode.FIRST_LOVE,
		value=rows[-1][1],
		ingest_id=rows[-1][0].ingest_id,
	)
	next_rows = repository.select_first_love(cursor=cursor, limit=2)
	assert [row[0].ingest_id for row in next_rows] == [second_tie.ingest_id, third.ingest_id]


def test_select_hall_of_fame_orders_by_hall_of_fame_at(session: Session) -> None:
	now = datetime(2024, 1, 2, tzinfo=timezone.utc)
	first = add_image_record(session, 4, captured_at=now)
	second = add_image_record(session, 3, captured_at=now - timedelta(hours=1))
	second_tie = add_image_record(session, 2, captured_at=now - timedelta(hours=2))
	third = add_image_record(session, 1, captured_at=now - timedelta(hours=3))

	add_stats_record(
		session,
		first.ingest_id,
		view_count=1,
		hall_of_fame_at=now,
	)
	add_stats_record(
		session,
		second.ingest_id,
		view_count=1,
		hall_of_fame_at=now - timedelta(days=2),
	)
	add_stats_record(
		session,
		second_tie.ingest_id,
		view_count=1,
		hall_of_fame_at=now - timedelta(days=2),
	)
	add_stats_record(
		session,
		third.ingest_id,
		view_count=1,
		hall_of_fame_at=now - timedelta(days=3),
	)

	repository = BaseImageListRepository(session, engaged_score_threshold=160)

	rows = repository.select_hall_of_fame(cursor=None, limit=2)
	assert [row[0].ingest_id for row in rows] == [first.ingest_id, second.ingest_id]

	cursor = DatetimeImageListCursor(
		mode=ImageListCursorMode.HALL_OF_FAME,
		value=rows[-1][1],
		ingest_id=rows[-1][0].ingest_id,
	)
	next_rows = repository.select_hall_of_fame(cursor=cursor, limit=2)
	assert [row[0].ingest_id for row in next_rows] == [second_tie.ingest_id, third.ingest_id]


def test_select_engaged_orders_by_score_evaluated(session: Session) -> None:
	now = datetime(2024, 1, 2, tzinfo=timezone.utc)
	first = add_image_record(session, 5, ingested_at=now)
	second = add_image_record(session, 4, ingested_at=now - timedelta(hours=1))
	second_tie = add_image_record(session, 3, ingested_at=now - timedelta(hours=2))
	third = add_image_record(session, 2, ingested_at=now - timedelta(hours=3))
	fourth = add_image_record(session, 1, ingested_at=now - timedelta(hours=4))
	fifth = add_image_record(session, 6, ingested_at=now - timedelta(hours=5))

	add_stats_record(
		session,
		first.ingest_id,
		view_count=1,
		score_evaluated=180,
		score_evaluated_at=now,
	)
	add_stats_record(
		session,
		second.ingest_id,
		view_count=1,
		score_evaluated=170,
		score_evaluated_at=now,
	)
	add_stats_record(
		session,
		second_tie.ingest_id,
		view_count=1,
		score_evaluated=170,
		score_evaluated_at=now,
	)
	add_stats_record(
		session,
		third.ingest_id,
		view_count=1,
		score_evaluated=165,
		score_evaluated_at=now,
	)
	add_stats_record(
		session,
		fourth.ingest_id,
		view_count=1,
		score_evaluated=150,
		score_evaluated_at=now,
	)
	add_stats_record(
		session,
		fifth.ingest_id,
		view_count=1,
		score_evaluated=190,
		score_evaluated_at=now,
		hall_of_fame_at=now - timedelta(hours=2),
	)
	session.commit()

	repository = BaseImageListRepository(session, engaged_score_threshold=160)

	rows = repository.select_engaged(cursor=None, limit=2)
	assert [row[0].ingest_id for row in rows] == [first.ingest_id, second.ingest_id]

	cursor = UInt8ImageListCursor(
		mode=ImageListCursorMode.ENGAGED,
		value=rows[-1][1],
		ingest_id=rows[-1][0].ingest_id,
	)
	assert cursor.value == 170

	next_rows = repository.select_engaged(cursor=cursor, limit=2)
	assert [row[0].ingest_id for row in next_rows] == [second_tie.ingest_id, third.ingest_id]
