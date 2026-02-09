from datetime import datetime, time, timezone
from types import SimpleNamespace
from zoneinfo import ZoneInfo

import pytest

from tests.stubs.user import StubUserRepository

import app.services.users.query_service as query_service
from app.domain.activities.daily_period import DailyPeriodResolver, resolve_daily_period_start
from app.models.records import UserRecord
from app.services.users.query_service import UserQueryService


def test_get_quota_returns_remaining_and_reset_at(monkeypatch: pytest.MonkeyPatch) -> None:
	current = datetime(2026, 1, 1, 12, 0, tzinfo=timezone.utc)
	monkeypatch.setattr(query_service, 'datetime', SimpleNamespace(now=lambda _: current))  # pyright: ignore[reportUnknownLambdaType]

	user = UserRecord(id=1)
	user.daily_love_used = 3
	user_repo = StubUserRepository()
	user_repo.users[user.id] = user

	service = UserQueryService(
		repository=user_repo,
		daily_love_limit=10,
		period_resolver=DailyPeriodResolver(
			base_timezone=ZoneInfo('UTC'),
			daily_reset_at=time(5, 0),
		),
	)

	response = service.get_quota()
	assert user_repo.get_called_with == [1]
	assert response.love.period == 'daily'
	assert response.love.reset_at == resolve_daily_period_start(
		current,
		daily_reset_at=time(5, 0),
		base_timezone=ZoneInfo('UTC'),
	)
	assert response.love.limit == 10
	assert response.love.remaining == 7


def test_get_quota_uses_limit_when_no_love_used(monkeypatch: pytest.MonkeyPatch) -> None:
	current = datetime(2026, 1, 1, 12, 0, tzinfo=timezone.utc)
	monkeypatch.setattr(query_service, 'datetime', SimpleNamespace(now=lambda _: current))  # pyright: ignore[reportUnknownLambdaType]

	user_repo = StubUserRepository()
	user_repo.users[1] = UserRecord(id=1)

	service = UserQueryService(
		repository=user_repo,
		daily_love_limit=8,
		period_resolver=DailyPeriodResolver(
			base_timezone=ZoneInfo('UTC'),
			daily_reset_at=time(5, 0),
		),
	)

	response = service.get_quota()
	assert user_repo.get_called_with == [1]
	assert response.love.period == 'daily'
	assert response.love.limit == 8
	assert response.love.remaining == 8


def test_get_quota_clamps_remaining_to_zero(monkeypatch: pytest.MonkeyPatch) -> None:
	current = datetime(2026, 1, 1, 12, 0, tzinfo=timezone.utc)
	monkeypatch.setattr(query_service, 'datetime', SimpleNamespace(now=lambda _: current))  # pyright: ignore[reportUnknownLambdaType]

	user = UserRecord(id=1)
	user.daily_love_used = 99
	user_repo = StubUserRepository()
	user_repo.users[user.id] = user

	service = UserQueryService(
		repository=user_repo,
		daily_love_limit=5,
		period_resolver=DailyPeriodResolver(
			base_timezone=ZoneInfo('UTC'),
			daily_reset_at=time(5, 0),
		),
	)

	response = service.get_quota()
	assert user_repo.get_called_with == [1]
	assert response.love.period == 'daily'
	assert response.love.limit == 5
	assert response.love.remaining == 0
