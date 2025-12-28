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


def test_upsert_with_increment(session: Session) -> None:
	repo = SQLiteStatsRepository(session)
	image = add_image_record(session, 20)

	stats = repo.upsert_with_increment(image.ingest_id)
	assert stats.view_count == 1
	assert stats.last_viewed_at is not None

	stats_again = repo.upsert_with_increment(image.ingest_id)
	assert stats_again.view_count == 2
