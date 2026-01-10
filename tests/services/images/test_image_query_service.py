from datetime import datetime, timedelta, timezone
from typing import Any, Generator

import pytest
from sqlmodel import Session, SQLModel, create_engine

from tests.services.activities.stats.factory import add_stats_record
from tests.services.images.utils import add_image_record

from app.config.environments import env
from app.persist.images.base import BaseImageRepository
from app.services.images.query_service import ImageQueryService


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	SQLModel.metadata.create_all(engine)
	with Session(engine) as session:
		yield session


def test_get_latest_orders_desc_and_sets_cursor(session: Session) -> None:
	now = datetime.now(timezone.utc)
	first = add_image_record(session, 1, ingested_at=now)
	second = add_image_record(session, 2, ingested_at=now - timedelta(hours=1))
	third = add_image_record(session, 3, ingested_at=now - timedelta(hours=2))

	service = ImageQueryService(
		session=session,
		repository=BaseImageRepository(session),
		engaged_score_threshold=160,
		variant_layers=env.variant_layers,
	)

	response = service.get_latest(cursor=None, limit=2, exclude_formats=())
	assert response.cursor == second.ingested_at
	assert response.items is not None
	assert [item.id for item in response.items] == [first.ingest_id, second.ingest_id]

	next_response = service.get_latest(cursor=response.cursor, limit=2, exclude_formats=())
	assert next_response.cursor is None
	assert next_response.items is not None
	assert [item.id for item in next_response.items] == [third.ingest_id]


def test_get_chronological_orders_by_captured_at(session: Session) -> None:
	now = datetime(2024, 1, 2, tzinfo=timezone.utc)
	first = add_image_record(session, 1, captured_at=now)
	second = add_image_record(session, 2, captured_at=now - timedelta(hours=1))

	service = ImageQueryService(
		session=session,
		repository=BaseImageRepository(session),
		engaged_score_threshold=160,
		variant_layers=env.variant_layers,
	)

	response = service.get_chronological(cursor=None, limit=1, exclude_formats=())
	assert response.cursor == now
	assert response.items is not None
	assert [item.id for item in response.items] == [first.ingest_id]

	next_response = service.get_chronological(cursor=response.cursor, limit=1, exclude_formats=())
	assert next_response.cursor is None
	assert next_response.items is not None
	assert [item.id for item in next_response.items] == [second.ingest_id]


def test_get_recently_orders_by_last_viewed_at(session: Session) -> None:
	now = datetime(2024, 1, 2, tzinfo=timezone.utc)
	first = add_image_record(session, 1, ingested_at=now)
	second = add_image_record(session, 2, ingested_at=now - timedelta(hours=1))

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

	service = ImageQueryService(
		session=session,
		repository=BaseImageRepository(session),
		engaged_score_threshold=160,
		variant_layers=env.variant_layers,
	)

	response = service.get_recently(cursor=None, limit=1, exclude_formats=())
	assert response.cursor == now
	assert response.items is not None
	assert [item.id for item in response.items] == [first.ingest_id]

	next_response = service.get_recently(cursor=response.cursor, limit=1, exclude_formats=())
	assert next_response.cursor is None
	assert next_response.items is not None
	assert [item.id for item in next_response.items] == [second.ingest_id]


def test_get_first_love_orders_by_first_loved_at(session: Session) -> None:
	now = datetime(2024, 1, 2, tzinfo=timezone.utc)
	first = add_image_record(session, 1, ingested_at=now)
	second = add_image_record(session, 2, ingested_at=now - timedelta(hours=1))

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

	service = ImageQueryService(
		session=session,
		repository=BaseImageRepository(session),
		engaged_score_threshold=160,
		variant_layers=env.variant_layers,
	)

	response = service.get_first_love(cursor=None, limit=1, exclude_formats=())
	assert response.cursor == now
	assert response.items is not None
	assert [item.id for item in response.items] == [first.ingest_id]

	next_response = service.get_first_love(cursor=response.cursor, limit=1, exclude_formats=())
	assert next_response.cursor is None
	assert next_response.items is not None
	assert [item.id for item in next_response.items] == [second.ingest_id]


def test_get_hall_of_fame_orders_by_hall_of_fame_at(session: Session) -> None:
	now = datetime(2024, 1, 2, tzinfo=timezone.utc)
	first = add_image_record(session, 1, captured_at=now)
	second = add_image_record(session, 2, captured_at=now - timedelta(hours=1))

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

	service = ImageQueryService(
		session=session,
		repository=BaseImageRepository(session),
		engaged_score_threshold=160,
		variant_layers=env.variant_layers,
	)

	response = service.get_hall_of_fame(cursor=None, limit=1, exclude_formats=())
	assert response.cursor == now
	assert response.items is not None
	assert [item.id for item in response.items] == [first.ingest_id]

	next_response = service.get_hall_of_fame(cursor=response.cursor, limit=1, exclude_formats=())
	assert next_response.cursor is None
	assert next_response.items is not None
	assert [item.id for item in next_response.items] == [second.ingest_id]


def test_get_engaged_orders_by_score_evaluated(session: Session) -> None:
	now = datetime(2024, 1, 2, tzinfo=timezone.utc)
	first = add_image_record(session, 1, ingested_at=now)
	second = add_image_record(session, 2, ingested_at=now - timedelta(hours=1))
	third = add_image_record(session, 3, ingested_at=now - timedelta(hours=2))
	fourth = add_image_record(session, 4, ingested_at=now - timedelta(hours=3))

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
		third.ingest_id,
		view_count=1,
		score_evaluated=150,
		score_evaluated_at=now,
	)
	add_stats_record(
		session,
		fourth.ingest_id,
		view_count=1,
		score_evaluated=190,
		score_evaluated_at=now,
		hall_of_fame_at=now - timedelta(hours=1, minutes=30),
	)
	session.commit()

	service = ImageQueryService(
		session=session,
		repository=BaseImageRepository(session),
		engaged_score_threshold=160,
		variant_layers=env.variant_layers,
	)

	response = service.get_engaged(cursor=None, limit=1, exclude_formats=())
	assert response.cursor == 180
	assert response.items is not None
	assert [item.id for item in response.items] == [first.ingest_id]

	next_response = service.get_engaged(cursor=response.cursor, limit=1, exclude_formats=())
	assert next_response.cursor is None
	assert next_response.items is not None
	assert [item.id for item in next_response.items] == [second.ingest_id]
