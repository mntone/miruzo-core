from datetime import datetime, timedelta, timezone
from typing import Any, Generator

import pytest
from sqlmodel import Session, SQLModel, create_engine

from tests.services.images.utils import add_image_record

from app.services.images.repository import ImageRepository


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	SQLModel.metadata.create_all(engine)
	with Session(engine) as session:
		yield session


def test_select_by_ingest_id(session: Session) -> None:
	repo = ImageRepository(session)
	now = datetime.now(timezone.utc)
	image = add_image_record(session, 10, captured_at=now)

	item = repo.select_by_ingest_id(image.ingest_id)
	assert item is not None
	assert item.ingest_id == image.ingest_id
	assert item.captured_at == image.captured_at
