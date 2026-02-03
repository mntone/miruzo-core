from datetime import datetime, time, timezone
from typing import final
from zoneinfo import ZoneInfo

import pytest

from tests.services.activities.stats.factory import build_stats_record
from tests.stubs.score import StubScoreCalculator
from tests.stubs.session import StubSession
from tests.stubs.stats import StubStatsRepository, create_stub_stats_repository_factory
from tests.stubs.user import StubUserRepository

from app.domain.activities.daily_period import DailyPeriodResolver
from app.models.enums import ActionKind
from app.models.records import ActionRecord
from app.persist.users.protocol import UserRepository
from app.services.activities.daily_decay import DailyDecayRunner


@final
class _StubDecayCreator:
	def __init__(self) -> None:
		self.called_with: list[tuple[int, datetime]] = []

	def create(self, ingest_id: int, *, occurred_at: datetime) -> ActionRecord:
		self.called_with.append((ingest_id, occurred_at))
		return ActionRecord(
			ingest_id=ingest_id,
			kind=ActionKind.DECAY,
			occurred_at=occurred_at,
		)


def test_apply_daily_decay_updates_scores(monkeypatch: pytest.MonkeyPatch) -> None:
	evaluated_at = datetime(2026, 1, 2, 6, 0, tzinfo=timezone.utc)
	stats_one = build_stats_record(1, score=10)
	stats_two = build_stats_record(
		2,
		score=20,
		last_viewed_at=datetime(2026, 1, 2, 5, 30, tzinfo=timezone.utc),
	)

	stats_repo = StubStatsRepository()
	stats_repo.stats_list_response = [stats_one, stats_two]

	decay_creator = _StubDecayCreator()
	score_calculator = StubScoreCalculator()
	user_repo = StubUserRepository()
	user = user_repo.get_or_create_singleton()
	user.daily_love_used = 3

	def _create_user_repository(_: StubSession) -> StubUserRepository:
		return user_repo

	class _DecayCreatorFactory:
		def __init__(self, *args, **kwargs) -> None:  # pyright: ignore[reportMissingParameterType, reportUnknownParameterType]
			pass

		def create(self, ingest_id: int, *, occurred_at: datetime) -> ActionRecord:
			return decay_creator.create(ingest_id, occurred_at=occurred_at)

	monkeypatch.setattr(
		'app.services.activities.daily_decay.create_stats_repository',
		create_stub_stats_repository_factory(stats_repo),
	)
	monkeypatch.setattr(
		'app.services.activities.daily_decay.create_user_repository',
		_create_user_repository,
	)
	monkeypatch.setattr(
		'app.services.activities.daily_decay.DecayActionCreator',
		_DecayCreatorFactory,
	)

	session = StubSession()
	runner = DailyDecayRunner(
		period_resolver=DailyPeriodResolver(
			base_timezone=ZoneInfo('UTC'),
			daily_reset_at=time(5, 0),
		),
		score_calculator=score_calculator,  # pyright: ignore[reportArgumentType]
	)

	runner.apply_daily_decay(
		session,  # pyright: ignore[reportArgumentType]
		evaluated_at=evaluated_at,
	)

	assert stats_one.score == 8
	assert stats_two.score == 18
	assert decay_creator.called_with == [
		(1, evaluated_at),
		(2, evaluated_at),
	]

	first_action, first_score, first_context = score_calculator.apply_called_with[0]
	assert first_action.kind == ActionKind.DECAY
	assert first_score == 10
	assert first_context.evaluated_at == evaluated_at
	assert first_context.has_view_today is False

	second_action, second_score, second_context = score_calculator.apply_called_with[1]
	assert second_action.kind == ActionKind.DECAY
	assert second_score == 20
	assert second_context.evaluated_at == evaluated_at
	assert second_context.has_view_today is True

	assert user.daily_love_used == 0


def test_apply_daily_decay_skips_when_no_action(monkeypatch: pytest.MonkeyPatch) -> None:
	evaluated_at = datetime(2026, 1, 2, 6, 0, tzinfo=timezone.utc)
	stats_one = build_stats_record(1, score=10)
	stats_two = build_stats_record(2, score=20)

	stats_repo = StubStatsRepository()
	stats_repo.stats_list_response = [stats_one, stats_two]

	score_calculator = StubScoreCalculator()
	user_repo = StubUserRepository()
	user = user_repo.get_or_create_singleton()
	user.daily_love_used = 2

	def _create_user_repository(_: StubSession) -> UserRepository:
		return user_repo

	class _DecayCreatorFactory:
		def __init__(self, *args, **kwargs) -> None:  # pyright: ignore[reportMissingParameterType, reportUnknownParameterType]
			pass

		def create(self, ingest_id: int, *, occurred_at: datetime) -> ActionRecord | None:
			if ingest_id == 2:
				return None
			return ActionRecord(
				ingest_id=ingest_id,
				kind=ActionKind.DECAY,
				occurred_at=occurred_at,
			)

	monkeypatch.setattr(
		'app.services.activities.daily_decay.create_stats_repository',
		create_stub_stats_repository_factory(stats_repo),
	)
	monkeypatch.setattr(
		'app.services.activities.daily_decay.create_user_repository',
		_create_user_repository,
	)
	monkeypatch.setattr(
		'app.services.activities.daily_decay.DecayActionCreator',
		_DecayCreatorFactory,
	)

	session = StubSession()
	runner = DailyDecayRunner(
		period_resolver=DailyPeriodResolver(
			base_timezone=ZoneInfo('UTC'),
			daily_reset_at=time(5, 0),
		),
		score_calculator=score_calculator,  # pyright: ignore[reportArgumentType]
	)

	runner.apply_daily_decay(
		session,  # pyright: ignore[reportArgumentType]
		evaluated_at=evaluated_at,
	)

	assert stats_one.score == 8
	assert stats_two.score == 20
	assert len(score_calculator.apply_called_with) == 1

	assert user.daily_love_used == 0
