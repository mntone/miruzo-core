from datetime import datetime, time, timezone
from zoneinfo import ZoneInfo

from tests.services.activities.stats.factory import build_stats_record

from app.config.score import ScoreConfig
from app.domain.activities.daily_period import DailyPeriodResolver
from app.domain.score.calculator import ScoreCalculator
from app.models.enums import ActionKind
from app.models.records import ActionRecord
from app.services.activities.stats.score_updater import update_score_from_action


def _make_resolver() -> DailyPeriodResolver:
	return DailyPeriodResolver(
		base_timezone=ZoneInfo('UTC'),
		daily_reset_at=time(0, 0),
	)


def _make_action(kind: ActionKind, occurred_at: datetime) -> ActionRecord:
	return ActionRecord(
		ingest_id=1,
		kind=kind,
		occurred_at=occurred_at,
	)


def _expected_score(config: ScoreConfig, current: int) -> int:
	return max(
		config.minimum_score,
		min(config.maximum_score, current + config.love_bonus),
	)


def test_update_score_from_action_updates_score_only() -> None:
	config = ScoreConfig()
	evaluated_at = datetime(2024, 1, 2, tzinfo=timezone.utc)
	stats = build_stats_record(
		1,
		score=100,
		score_evaluated=80,
		score_evaluated_at=datetime(2024, 1, 1, tzinfo=timezone.utc),
	)

	update_score_from_action(
		stats=stats,
		action=_make_action(ActionKind.LOVE, evaluated_at),
		evaluated_at=evaluated_at,
		resolver=_make_resolver(),
		score_calculator=ScoreCalculator(config),
	)

	assert stats.score == _expected_score(config, 100)
	assert stats.score_evaluated == 80
	assert stats.score_evaluated_at == datetime(2024, 1, 1, tzinfo=timezone.utc)


def test_update_score_from_action_updates_evaluated_fields() -> None:
	config = ScoreConfig()
	evaluated_at = datetime(2024, 1, 2, tzinfo=timezone.utc)
	stats = build_stats_record(
		1,
		score=150,
		score_evaluated=120,
		score_evaluated_at=None,
	)

	update_score_from_action(
		stats=stats,
		action=_make_action(ActionKind.LOVE, evaluated_at),
		evaluated_at=evaluated_at,
		resolver=_make_resolver(),
		score_calculator=ScoreCalculator(config),
		update_evaluated=True,
	)

	assert stats.score == _expected_score(config, 150)
	assert stats.score_evaluated == stats.score
	assert stats.score_evaluated_at == evaluated_at
