from app.config.environments import env
from app.models.api.activities.stats import LoveStatsModel, StatsModel
from app.models.records import StatsRecord


def test_stats_model_from_record_clamps_score_to_public_minimum() -> None:
	stats = StatsRecord(
		ingest_id=1,
		score=-10,
		score_evaluated=100,
		view_count=0,
		last_viewed_at=None,
	)

	model = StatsModel.from_record(stats)

	assert model.score == env.score.public_minimum_score


def test_love_stats_model_from_record_clamps_score_to_public_minimum() -> None:
	stats = StatsRecord(
		ingest_id=1,
		score=-10,
		score_evaluated=100,
		view_count=0,
		last_viewed_at=None,
	)

	love_model = LoveStatsModel.from_record(stats)

	assert love_model.score == env.score.public_minimum_score
