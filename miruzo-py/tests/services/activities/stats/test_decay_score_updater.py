from datetime import datetime, timedelta, timezone

from tests.services.activities.stats.factory import build_stats_record
from tests.stubs.decay_score import StubDecayScoreCalculator

from app.domain.activities.daily_period import DailyPeriodResolver
from app.services.activities.stats.decay_score_updater import update_decay_score


def _make_resolver() -> DailyPeriodResolver:
	return DailyPeriodResolver(timedelta())


def _expected_score(current: int, delta: int) -> int:
	return current + delta


def test_update_decay_score_updates_score_only() -> None:
	score_calculator = StubDecayScoreCalculator()
	evaluated_at = datetime(2024, 1, 2, tzinfo=timezone.utc)
	stats = build_stats_record(
		1,
		score=100,
		score_evaluated=80,
		score_evaluated_at=datetime(2024, 1, 1, tzinfo=timezone.utc),
	)

	update_decay_score(
		stats=stats,
		evaluated_at=evaluated_at,
		resolver=_make_resolver(),
		decay_score_calculator=score_calculator,  # pyright: ignore[reportArgumentType]
	)

	assert stats.score == _expected_score(100, -2)
	assert stats.score_evaluated == _expected_score(100, -2)
	assert stats.score_evaluated_at == evaluated_at


def test_update_decay_score_updates_evaluated_fields() -> None:
	score_calculator = StubDecayScoreCalculator()
	evaluated_at = datetime(2024, 1, 2, tzinfo=timezone.utc)
	stats = build_stats_record(
		1,
		score=150,
		score_evaluated=120,
		score_evaluated_at=None,
	)

	update_decay_score(
		stats=stats,
		evaluated_at=evaluated_at,
		resolver=_make_resolver(),
		decay_score_calculator=score_calculator,  # pyright: ignore[reportArgumentType]
	)

	assert stats.score == _expected_score(150, -2)
	assert stats.score_evaluated == _expected_score(150, -2)
	assert stats.score_evaluated_at == evaluated_at
