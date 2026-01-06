from datetime import datetime, time, timezone
from typing import Any, Generator
from zoneinfo import ZoneInfo

import pytest
from sqlmodel import Session, SQLModel, create_engine

from tests.services.images.utils import add_ingest_record

from app.config.score import ScoreConfig
from app.errors import InvalidStateError
from app.models.enums import ActionKind
from app.services.activities.actions.repository import ActionRepository
from app.services.activities.love_cancel import LoveCancelRunner
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


def test_run_restores_previous_love(session: Session) -> None:
	action_repo = ActionRepository(session)
	stats_repo = SQLiteStatsRepository(session)
	user_repo = SQLiteUserRepository(session)

	period_start = datetime(2024, 1, 2, tzinfo=timezone.utc)
	previous_love = datetime(2024, 1, 1, 23, 0, tzinfo=timezone.utc)
	with session.begin():
		user = user_repo.get_or_create_singleton()
		user.daily_love_used = 1

		ingest = add_ingest_record(session, 1)
		stats = stats_repo.get_or_create(ingest.id, initial_score=100)
		stats.first_loved_at = previous_love
		stats.last_loved_at = datetime(2024, 1, 2, 1, 0, tzinfo=timezone.utc)

		action_repo.insert(
			ingest.id,
			kind=ActionKind.LOVE,
			occurred_at=previous_love,
		)

	runner = LoveCancelRunner(
		base_timezone=ZoneInfo('UTC'),
		daily_reset_at=time(0, 0),
		score_config=ScoreConfig(),
	)

	with session.begin():
		response = runner.run(session, ingest_id=ingest.id, evaluated_at=period_start)
	assert response.stats.first_loved_at == previous_love
	assert response.stats.last_loved_at == previous_love
	assert response.stats.score == 82

	user = user_repo.get_or_create_singleton()
	assert user.daily_love_used == 0

	stats = stats_repo.get_one(ingest.id)
	assert stats.first_loved_at == previous_love
	assert stats.last_loved_at == previous_love
	assert stats.score == 82


def test_run_clears_first_loved_at_when_no_previous_love(session: Session) -> None:
	stats_repo = SQLiteStatsRepository(session)

	evaluated_at = datetime(2024, 1, 2, 1, 0, tzinfo=timezone.utc)
	with session.begin():
		ingest = add_ingest_record(session, 1)
		stats = stats_repo.get_or_create(ingest.id, initial_score=100)
		stats.first_loved_at = evaluated_at
		stats.last_loved_at = evaluated_at

	runner = LoveCancelRunner(
		base_timezone=ZoneInfo('UTC'),
		daily_reset_at=time(0, 0),
		score_config=ScoreConfig(),
	)

	with session.begin():
		response = runner.run(session, ingest_id=ingest.id, evaluated_at=evaluated_at)
	assert response.stats.first_loved_at is None
	assert response.stats.last_loved_at is None

	stats = stats_repo.get_one(ingest.id)
	assert stats.first_loved_at is None
	assert stats.last_loved_at is None


def test_run_raises_when_no_love_in_period(session: Session) -> None:
	stats_repo = SQLiteStatsRepository(session)

	with session.begin():
		ingest = add_ingest_record(session, 1)
		stats = stats_repo.get_or_create(ingest.id, initial_score=100)
		stats.last_loved_at = datetime(2024, 1, 1, 23, 0, tzinfo=timezone.utc)

	runner = LoveCancelRunner(
		base_timezone=ZoneInfo('UTC'),
		daily_reset_at=time(0, 0),
		score_config=ScoreConfig(),
	)

	with pytest.raises(InvalidStateError):
		with session.begin():
			runner.run(
				session,
				ingest_id=ingest.id,
				evaluated_at=datetime(2024, 1, 2, tzinfo=timezone.utc),
			)


def test_run_raises_when_first_loved_at_missing(session: Session) -> None:
	stats_repo = SQLiteStatsRepository(session)

	evaluated_at = datetime(2024, 1, 2, 1, 0, tzinfo=timezone.utc)
	with session.begin():
		ingest = add_ingest_record(session, 1)
		stats = stats_repo.get_or_create(ingest.id, initial_score=100)
		stats.first_loved_at = None
		stats.last_loved_at = datetime(2024, 1, 2, 0, 30, tzinfo=timezone.utc)

	runner = LoveCancelRunner(
		base_timezone=ZoneInfo('UTC'),
		daily_reset_at=time(0, 0),
		score_config=ScoreConfig(),
	)

	with pytest.raises(InvalidStateError):
		with session.begin():
			runner.run(session, ingest_id=ingest.id, evaluated_at=evaluated_at)
