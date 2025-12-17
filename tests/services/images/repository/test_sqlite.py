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

	items, cursor = repo.get_list(cursor=None, limit=2)

	assert [item.id for item in items] == [first.id, second.id]
	assert cursor == items[-1].captured_at

	next_items, next_cursor = repo.get_list(cursor=cursor, limit=2)
	assert [item.id for item in next_items] == [third.id]
	assert next_cursor is None


def test_get_detail_with_stats(session: Session) -> None:
	repo = SQLiteImageRepository(session)
	image = add_image_record(session, 10, captured_at=datetime.now(timezone.utc))
	stats = StatsRecord(
		image_id=image.id,
		favorite=True,
		score=42,
		view_count=5,
		last_viewed_at=datetime.now(timezone.utc),
	)
	session.add(stats)
	session.commit()

	detail = repo.get_detail(image.id)
	assert detail is not None
	assert detail.id == image.id

	result = repo.get_detail_with_stats(image.id)
	assert result is not None
	image_record, stats_record = result
	assert image_record.id == image.id
	assert stats_record is not None
	assert stats_record.favorite is True
	assert stats_record.score == 42


def test_get_stats_with_increment_creates_and_updates(session: Session) -> None:
	repo = SQLiteImageRepository(session)
	image = add_image_record(session, 20, captured_at=datetime.now(timezone.utc))

	stats = repo.upsert_stats_with_increment(image.id)
	assert stats.view_count == 1
	assert stats.last_viewed_at is not None

	stats_again = repo.upsert_stats_with_increment(image.id)
	assert stats_again.view_count == 2


def test_update_favorite_and_score(session: Session) -> None:
	repo = SQLiteImageRepository(session)
	image = add_image_record(session, 30, captured_at=datetime.now(timezone.utc))
	repo.create_stats(image.id)

	favorite_response = repo.update_favorite(image.id, True)
	assert favorite_response is not None
	stats_record = session.get(StatsRecord, image.id)
	assert stats_record is not None
	assert stats_record.favorite is True

	score_response = repo.update_score(image.id, 10)
	assert score_response is not None
	stats_record = session.get(StatsRecord, image.id)
	assert stats_record.score == DEFAULT_SCORE + 10
