from datetime import datetime, timedelta, timezone
from typing import Any, Generator

import pytest
from sqlmodel import Session, SQLModel, create_engine

from tests.services.images.utils import add_image_record

from app.config.constants import DEFAULT_SCORE
from app.models.records import StatsRecord
from app.services.images.repository.sqlite import SQLiteImageRepository


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	SQLModel.metadata.create_all(engine)
	with Session(engine) as session:
		yield session


def test_get_list_orders_desc_and_returns_cursor(session: Session) -> None:
	repo = SQLiteImageRepository(session)
	now = datetime.now(timezone.utc)
	first = add_image_record(session, 1, captured_at=now)
	second = add_image_record(session, 2, captured_at=now - timedelta(hours=1))
	third = add_image_record(session, 3, captured_at=now - timedelta(hours=2))

	items, cursor = repo.get_latest(cursor=None, limit=2)

	assert [item.ingest_id for item in items] == [first.ingest_id, second.ingest_id]
	assert cursor == items[-1].captured_at

	next_items, next_cursor = repo.get_latest(cursor=cursor, limit=2)
	assert [item.ingest_id for item in next_items] == [third.ingest_id]
	assert next_cursor is None


def test_get_detail_with_stats(session: Session) -> None:
	repo = SQLiteImageRepository(session)
	image = add_image_record(session, 10)
	current = datetime.now(timezone.utc)
	stats = StatsRecord(
		ingest_id=image.ingest_id,
		hall_of_fame_at=current,
		score=42,
		view_count=5,
		last_viewed_at=current,
	)
	session.add(stats)
	session.commit()

	detail = repo.get_context(image.ingest_id)
	assert detail is not None
	assert detail.ingest_id == image.ingest_id

	result = repo.get_detail_with_stats(image.ingest_id)
	assert result is not None
	image_record, stats_record = result
	assert image_record.ingest_id == image.ingest_id
	assert stats_record is not None
	assert stats_record.hall_of_fame_at == current.replace(tzinfo=None)
	assert stats_record.score == 42
	assert stats_record.view_count == 5


def test_upsert_stats_with_increment(session: Session) -> None:
	repo = SQLiteImageRepository(session)
	image = add_image_record(session, 20)

	stats = repo.upsert_stats_with_increment(image.ingest_id)
	assert stats.view_count == 1
	assert stats.last_viewed_at is not None

	stats_again = repo.upsert_stats_with_increment(image.ingest_id)
	assert stats_again.view_count == 2


def test_update_hall_of_fame(session: Session) -> None:
	repo = SQLiteImageRepository(session)
	image = add_image_record(session, 30)
	repo.create_stats(image.ingest_id)

	favorite_response = repo.update_favorite(image.ingest_id, True)
	assert favorite_response is not None
	assert favorite_response.favorited_at is not None

	stats = session.get(StatsRecord, image.ingest_id)
	assert stats is not None
	assert stats.hall_of_fame_at == favorite_response.favorited_at


def test_update_score(session: Session) -> None:
	repo = SQLiteImageRepository(session)
	image = add_image_record(session, 40)
	repo.create_stats(image.ingest_id)

	score_response = repo.update_score(image.ingest_id, 5)
	assert score_response is not None
	assert score_response.score == DEFAULT_SCORE + 5

	stats = session.get(StatsRecord, image.ingest_id)
	assert stats is not None
	assert stats.score == score_response.score
