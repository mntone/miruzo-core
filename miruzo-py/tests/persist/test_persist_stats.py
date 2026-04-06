from datetime import datetime, timezone
from typing import Any, Generator

import pytest
from sqlalchemy import create_engine
from sqlalchemy.orm import Session

from tests.persist.utils import add_ingest_row, get_stats_row

from app.databases.metadata import metadata
from app.persist.stats.implementation import create_stats_repository
from app.persist.stats.protocol import StatsCreateInput


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	metadata.create_all(engine)
	with Session(engine) as session:
		yield session


def test_create_persists_stats_row(session: Session) -> None:
	now = datetime(2026, 1, 1, tzinfo=timezone.utc)
	ingest_id = add_ingest_row(session, ingested_at=now)

	create_stats_repository(session).create(
		StatsCreateInput(
			ingest_id=ingest_id,
			initial_score=42,
		),
	)

	row = get_stats_row(session, ingest_id=ingest_id)
	assert row['ingest_id'] == ingest_id
	assert row['score'] == 42
	assert row['score_evaluated'] == 42
