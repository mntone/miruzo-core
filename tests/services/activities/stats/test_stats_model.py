from app.config.environments import env
from app.models.api.activities.stats import StatsModel
from app.models.records import StatsRecord


def test_from_record_clamps_score_to_minimum() -> None:
	stats = StatsRecord(
		ingest_id=1,
		score=-10,
		score_evaluated=100,
		view_count=0,
		last_viewed_at=None,
	)

	model = StatsModel.from_record(stats)

	assert model.score == env.score.minimum_score
