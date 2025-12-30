from collections.abc import Generator
from datetime import datetime, timezone
from typing import Any

import pytest
from sqlmodel import Session, SQLModel, create_engine

from tests.services.images.utils import add_ingest_record

from app.models.enums import ActionKind
from app.services.activities.actions.repository import ActionRepository


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	SQLModel.metadata.create_all(engine)
	with Session(engine) as session:
		yield session


def test_select_by_ingest_id_orders_by_occurred_at(session: Session) -> None:
	ingest = add_ingest_record(session, 1)
	repo = ActionRepository(session)

	repo.insert(
		ingest.id,
		kind=ActionKind.VIEW,
		occurred_at=datetime(2024, 1, 2, tzinfo=timezone.utc),
	)
	repo.insert(
		ingest.id,
		kind=ActionKind.MEMO,
		occurred_at=datetime(2024, 1, 1, tzinfo=timezone.utc),
	)

	items = repo.select_by_ingest_id(ingest.id)

	assert [item.kind for item in items] == [
		ActionKind.MEMO,
		ActionKind.VIEW,
	]
