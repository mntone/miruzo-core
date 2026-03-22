from app.config.score import ScoreConfig
from app.domain.decay_score.calculator import DecayScoreCalculator
from app.domain.decay_score.context import DecayScoreContext


def _make_context(
	*,
	days_since_last_view: int,
	has_view_today: bool,
) -> DecayScoreContext:
	return DecayScoreContext(
		days_since_last_view=days_since_last_view,
		has_view_today=has_view_today,
	)


def test_apply_decay_high_score_with_no_access() -> None:
	calc = DecayScoreCalculator(ScoreConfig())
	context = _make_context(
		days_since_last_view=5,
		has_view_today=False,
	)

	score = calc.apply(
		score=200,
		context=context,
	)

	assert score == 198


def test_apply_decay_interval10d_penalty() -> None:
	calc = DecayScoreCalculator(ScoreConfig())
	context = _make_context(
		days_since_last_view=10,
		has_view_today=True,
	)

	score = calc.apply(
		score=100,
		context=context,
	)

	assert score == 97


def test_apply_clamps_score() -> None:
	config = ScoreConfig(
		minimum_score=10,
		maximum_score=20,
	)
	calc = DecayScoreCalculator(config)
	context = _make_context(
		days_since_last_view=0,
		has_view_today=False,
	)

	assert calc.apply(score=50, context=context) == 20
	assert calc.apply(score=0, context=context) == 10
