from datetime import datetime, time, timezone

from tests.services.activities.actions.stubs import StubActionRepository

from app.domain.activities.daily_period import resolve_daily_period_range
from app.models.enums import ActionKind
from app.models.records import ActionRecord
from app.services.activities.actions.decay_creator import DecayActionCreator


def test_create_returns_action_when_missing() -> None:
	evaluated_at = datetime(2026, 1, 2, 6, 0, tzinfo=timezone.utc)
	repo = StubActionRepository()

	creator = DecayActionCreator(
		repository=repo,  # pyright: ignore[reportArgumentType]
		daily_reset_at=time(5, 0),
	)

	result = creator.create(1, occurred_at=evaluated_at)
	assert result is not None
	assert len(repo.insert_called_with) == 1
	assert repo.insert_called_with[0].ingest_id == result.ingest_id
	assert repo.insert_called_with[0].kind == result.kind
	assert repo.insert_called_with[0].occurred_at == result.occurred_at

	expected_since, expected_until = resolve_daily_period_range(
		evaluated_at,
		daily_reset_at=time(5, 0),
	)
	assert repo.select_one_called_with is not None
	assert repo.select_one_called_with.ingest_id == 1
	assert repo.select_one_called_with.kind == ActionKind.DECAY
	assert repo.select_one_called_with.since_occurred_at == expected_since
	assert repo.select_one_called_with.until_occurred_at == expected_until


def test_create_returns_none_when_existing() -> None:
	evaluated_at = datetime(2026, 1, 2, 6, 0, tzinfo=timezone.utc)
	existing = ActionRecord(
		ingest_id=1,
		kind=ActionKind.DECAY,
		occurred_at=evaluated_at,
	)

	repo = StubActionRepository()
	repo.actions = [existing]

	creator = DecayActionCreator(
		repository=repo,  # pyright: ignore[reportArgumentType]
		daily_reset_at=time(5, 0),
	)

	result = creator.create(1, occurred_at=evaluated_at)
	assert result is None
	assert len(repo.insert_called_with) == 0
