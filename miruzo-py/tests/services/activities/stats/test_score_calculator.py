from datetime import datetime, timezone

from app.config.score import ScoreConfig
from app.domain.score.calculator import ScoreCalculator
from app.domain.score.context import ScoreContext
from app.models.enums import ActionKind
from app.models.records import ActionRecord


def _make_action(kind: ActionKind) -> ActionRecord:
	return ActionRecord(
		ingest_id=1,
		kind=kind,
		occurred_at=datetime(2024, 1, 1, tzinfo=timezone.utc),
	)


def _make_context(
	*,
	last_viewed_at: datetime | None,
	days_since_last_view: int,
	has_view_today: bool,
) -> ScoreContext:
	return ScoreContext(
		evaluated_at=datetime(2024, 1, 2, tzinfo=timezone.utc),
		last_viewed_at=last_viewed_at,
		days_since_last_view=days_since_last_view,
		has_view_today=has_view_today,
	)


def test_apply_view_first_bonus() -> None:
	calc = ScoreCalculator(ScoreConfig())
	context = _make_context(
		last_viewed_at=None,
		days_since_last_view=0,
		has_view_today=False,
	)

	score = calc.apply(
		action=_make_action(ActionKind.VIEW),
		score=50,
		context=context,
	)

	assert score == 60


def test_apply_view_bonus_by_days() -> None:
	calc = ScoreCalculator(ScoreConfig())
	context = _make_context(
		last_viewed_at=datetime(2024, 1, 1, tzinfo=timezone.utc),
		days_since_last_view=2,
		has_view_today=False,
	)

	score = calc.apply(
		action=_make_action(ActionKind.VIEW),
		score=50,
		context=context,
	)

	assert score == 57


def test_apply_decay_high_score_with_no_access() -> None:
	calc = ScoreCalculator(ScoreConfig())
	context = _make_context(
		last_viewed_at=datetime(2024, 1, 1, tzinfo=timezone.utc),
		days_since_last_view=5,
		has_view_today=False,
	)

	score = calc.apply(
		action=_make_action(ActionKind.DECAY),
		score=200,
		context=context,
	)

	assert score == 198


def test_apply_decay_interval10d_penalty() -> None:
	calc = ScoreCalculator(ScoreConfig())
	context = _make_context(
		last_viewed_at=datetime(2024, 1, 1, tzinfo=timezone.utc),
		days_since_last_view=10,
		has_view_today=False,
	)

	score = calc.apply(
		action=_make_action(ActionKind.DECAY),
		score=100,
		context=context,
	)

	assert score == 98


def test_apply_clamps_score() -> None:
	config = ScoreConfig(
		minimum_score=10,
		maximum_score=20,
	)
	calc = ScoreCalculator(config)
	context = _make_context(
		last_viewed_at=None,
		days_since_last_view=0,
		has_view_today=False,
	)
	action = _make_action(ActionKind.UNKNOWN)

	assert calc.apply(action=action, score=50, context=context) == 20
	assert calc.apply(action=action, score=0, context=context) == 10
