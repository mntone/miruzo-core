from typing import Any, Generator

import pytest
from sqlmodel import Session, create_engine

from tests.persist.utils import get_stats_row
from tests.services.images.utils import add_ingest_record

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
	ingest = add_ingest_record(session, 20)
	create_stats_repository(session).create(
		StatsCreateInput(
			ingest_id=ingest.id,
			initial_score=42,
		),
	)

	row = get_stats_row(session, ingest_id=ingest.id)
	assert row['ingest_id'] == ingest.id
	assert row['score'] == 42
	assert row['score_evaluated'] == 42
