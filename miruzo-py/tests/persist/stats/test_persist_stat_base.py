from typing import Any, Generator

import pytest
from sqlmodel import Session, SQLModel, create_engine

from tests.services.images.utils import add_image_record

from app.persist.stats.base import BaseStatsRepository


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	SQLModel.metadata.create_all(engine)
	with Session(engine) as session:
		yield session


def test_create(session: Session) -> None:
	repo = BaseStatsRepository(session)
	image = add_image_record(session, 20)

	stats = repo.create(image.ingest_id, initial_score=42)
	assert stats.score == 42
	assert stats.view_count == 0
	assert stats.last_viewed_at is None
