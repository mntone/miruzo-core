from datetime import datetime
from typing import Any, Generator

import pytest
from sqlmodel import Session, SQLModel, create_engine

from app.services.images.query_executor import ImageListQueryExecutor


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	SQLModel.metadata.create_all(engine)
	with Session(engine) as session:
		yield session


def test_latest_adds_cursor_filter_and_order(session: Session) -> None:
	cursor = datetime(2024, 1, 1)
	statement = (
		ImageListQueryExecutor(session, engaged_score_threshold=160)
		.latest(cursor=cursor)
		.order_by_ingest_id()
		.limit(4)
		._statement
	)

	sql = str(statement)

	assert 'FROM images' in sql
	assert 'images.ingested_at <' in sql
	assert 'ORDER BY images.ingested_at DESC' in sql
	assert statement is not None
	assert statement._limit_clause is not None
	assert statement._limit_clause.value == 4


def test_chronological_filters_on_captured_at(session: Session) -> None:
	cursor = datetime(2024, 1, 1)
	statement = (
		ImageListQueryExecutor(session, engaged_score_threshold=160)
		.chronological(cursor=cursor)
		.order_by_ingest_id()
		.limit(5)
		._statement
	)

	sql = str(statement)

	assert 'JOIN ingests ON ingests.id = images.ingest_id' in sql
	assert 'ingests.captured_at <' in sql
	assert 'ORDER BY ingests.captured_at DESC' in sql
	assert statement is not None
	assert statement._limit_clause is not None
	assert statement._limit_clause.value == 5


def test_recently_filters_on_last_viewed_at(session: Session) -> None:
	cursor = datetime(2024, 1, 1)
	statement = (
		ImageListQueryExecutor(session, engaged_score_threshold=160)
		.recently(cursor=cursor)
		.order_by_ingest_id()
		.limit(6)
		._statement
	)

	sql = str(statement)

	assert 'JOIN stats ON stats.ingest_id = images.ingest_id' in sql
	assert 'stats.last_viewed_at IS NOT NULL' in sql
	assert 'stats.last_viewed_at <' in sql
	assert 'ORDER BY stats.last_viewed_at DESC' in sql
	assert statement is not None
	assert statement._limit_clause is not None
	assert statement._limit_clause.value == 6


def test_first_love_filters_on_first_loved_at(session: Session) -> None:
	cursor = datetime(2024, 1, 1)
	statement = (
		ImageListQueryExecutor(session, engaged_score_threshold=160)
		.first_love(cursor=cursor)
		.order_by_ingest_id()
		.limit(7)
		._statement
	)

	sql = str(statement)

	assert 'JOIN stats ON stats.ingest_id = images.ingest_id' in sql
	assert 'stats.first_loved_at IS NOT NULL' in sql
	assert 'stats.first_loved_at <' in sql
	assert 'ORDER BY stats.first_loved_at DESC' in sql
	assert statement is not None
	assert statement._limit_clause is not None
	assert statement._limit_clause.value == 7


def test_hall_of_fame_filters_on_hall_of_fame_at(session: Session) -> None:
	cursor = datetime(2024, 1, 1)
	statement = (
		ImageListQueryExecutor(session, engaged_score_threshold=160)
		.hall_of_fame(cursor=cursor)
		.order_by_ingest_id()
		.limit(8)
		._statement
	)

	sql = str(statement)

	assert 'JOIN stats ON stats.ingest_id = images.ingest_id' in sql
	assert 'stats.hall_of_fame_at IS NOT NULL' in sql
	assert 'stats.hall_of_fame_at <' in sql
	assert 'ORDER BY stats.hall_of_fame_at DESC' in sql
	assert statement is not None
	assert statement._limit_clause is not None
	assert statement._limit_clause.value == 8


def test_engaged_filters_on_score_evaluated(session: Session) -> None:
	cursor = 160
	statement = (
		ImageListQueryExecutor(session, engaged_score_threshold=160)
		.engaged(cursor=cursor)
		.order_by_ingest_id()
		.limit(9)
		._statement
	)

	sql = str(statement)

	assert 'JOIN stats ON stats.ingest_id = images.ingest_id' in sql
	assert 'stats.hall_of_fame_at IS NULL' in sql
	assert 'stats.score_evaluated >=' in sql
	assert 'stats.score_evaluated <' in sql
	assert 'ORDER BY stats.score_evaluated DESC' in sql
	assert statement is not None
	assert statement._limit_clause is not None
	assert statement._limit_clause.value == 9
