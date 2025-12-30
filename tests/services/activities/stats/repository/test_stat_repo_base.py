from typing import Any, Generator

import pytest
from sqlmodel import Session, SQLModel, create_engine

from tests.services.images.utils import add_image_record

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
