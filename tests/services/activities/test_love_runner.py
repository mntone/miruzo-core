from datetime import datetime, time, timezone
from typing import Any, Generator
from zoneinfo import ZoneInfo

import pytest
from sqlmodel import Session, SQLModel, create_engine

from tests.services.images.utils import add_ingest_record

from app.config.score import ScoreConfig
from app.errors import InvalidStateError, QuotaExceededError
from app.services.activities.love import LoveRunner
from app.services.activities.stats.repository.sqlite import SQLiteStatsRepository
from app.services.users.repository.sqlite import SQLiteUserRepository


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	SQLModel.metadata.create_all(engine)
	with Session(engine) as session:
		yield session


def test_run_updates_stats_and_user(session: Session) -> None:
	stats_repo = SQLiteStatsRepository(session)

	with session.begin():
		ingest = add_ingest_record(session, 1)
		stats_repo.get_or_create(ingest.id, initial_score=100)

	runner = LoveRunner(
		base_timezone=ZoneInfo('UTC'),
		daily_reset_at=time(0, 0),
		daily_love_limit=3,
		score_config=ScoreConfig(),
	)

	evaluated_at = datetime(2024, 1, 2, tzinfo=timezone.utc)
	with session.begin():
		response = runner.run(session, ingest_id=ingest.id, evaluated_at=evaluated_at)
	assert response.stats.first_loved_at == evaluated_at
	assert response.stats.last_loved_at == evaluated_at
	assert response.stats.score == 120

	user = SQLiteUserRepository(session).get_or_create_singleton()
	assert user.daily_love_used == 1

	stats = stats_repo.get_one(ingest.id)
	assert stats.first_loved_at == evaluated_at
	assert stats.last_loved_at == evaluated_at
	assert stats.score == 120


def test_run_raises_when_already_loved_today(session: Session) -> None:
	stats_repo = SQLiteStatsRepository(session)

	existing = datetime(2024, 1, 1, 1, tzinfo=timezone.utc)
	with session.begin():
		ingest = add_ingest_record(session, 1)

		stats = stats_repo.get_or_create(ingest.id, initial_score=100)
		stats.last_loved_at = existing

	runner = LoveRunner(
		base_timezone=ZoneInfo('UTC'),
		daily_reset_at=time(0, 0),
		daily_love_limit=3,
		score_config=ScoreConfig(),
	)

	with pytest.raises(InvalidStateError):
		with session.begin():
			runner.run(
				session,
				ingest_id=ingest.id,
				evaluated_at=datetime(2024, 1, 1, 2, tzinfo=timezone.utc),
			)

	stats = stats_repo.get_one(ingest.id)
	assert stats.last_loved_at == existing


def test_run_raises_when_quota_exceeded(session: Session) -> None:
	stats_repo = SQLiteStatsRepository(session)
	user_repo = SQLiteUserRepository(session)

	with session.begin():
		user = user_repo.get_or_create_singleton()
		user.daily_love_used = 1

		ingest = add_ingest_record(session, 1)
		stats_repo.get_or_create(ingest.id, initial_score=100)

	runner = LoveRunner(
		base_timezone=ZoneInfo('UTC'),
		daily_reset_at=time(0, 0),
		daily_love_limit=1,
		score_config=ScoreConfig(),
	)

	with pytest.raises(QuotaExceededError):
		with session.begin():
			runner.run(
				session,
				ingest_id=ingest.id,
				evaluated_at=datetime(2024, 1, 2, tzinfo=timezone.utc),
			)

	user = user_repo.get_or_create_singleton()
	assert user.daily_love_used == 1

	stats = stats_repo.get_one(ingest.id)
	assert stats.last_loved_at is None
