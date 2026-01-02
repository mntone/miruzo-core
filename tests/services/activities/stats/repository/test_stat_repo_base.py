from typing import Any, Generator

import pytest
from sqlmodel import Session, SQLModel, create_engine

from tests.services.images.utils import add_image_record, add_ingest_record

from app.services.activities.stats.repository.sqlite import SQLiteStatsRepository


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	SQLModel.metadata.create_all(engine)
	with Session(engine) as session:
		yield session


def test_get_or_create(session: Session) -> None:
	repo = SQLiteStatsRepository(session)
	image = add_image_record(session, 20)

	stats = repo.get_or_create(image.ingest_id, initial_score=42)
	assert stats.score == 42
	assert stats.view_count == 0
	assert stats.last_viewed_at is None

	stats_again = repo.get_or_create(image.ingest_id, initial_score=99)
	assert stats_again.ingest_id == image.ingest_id
	assert stats_again.score == 42


def test_iterable_paginates_without_duplicates(session: Session) -> None:
	repo = SQLiteStatsRepository(session)
	ingest_ids = [idx for idx in range(1, 503) if idx != 251]

	for ingest_id in ingest_ids:
		add_ingest_record(session, ingest_id)
		repo.get_or_create(ingest_id, initial_score=1)

	records = list(repo.iterable())
	record_ids = [record.ingest_id for record in records]

	assert record_ids == sorted(ingest_ids)
	assert len(record_ids) == len(set(record_ids))


def test_iterable_returns_empty_for_empty_table(session: Session) -> None:
	repo = SQLiteStatsRepository(session)

	assert list(repo.iterable()) == []
