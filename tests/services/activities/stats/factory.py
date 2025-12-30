from datetime import datetime

from app.models.records import StatsRecord


def _make_stats_record(
	ingest_id: int,
	*,
	score: int,
	view_count: int,
	last_viewed_at: datetime | None,
) -> StatsRecord:
	return StatsRecord(
		ingest_id=ingest_id,
		score=score,
		view_count=view_count,
		last_viewed_at=last_viewed_at,
		first_loved_at=None,
		hall_of_fame_at=None,
	)


def build_stats_record(
	ingest_id: int,
	*,
	score: int = 100,
	view_count: int = 0,
	last_viewed_at: datetime | None = None,
) -> StatsRecord:
	return _make_stats_record(
		ingest_id=ingest_id,
		score=score,
		view_count=view_count,
		last_viewed_at=last_viewed_at,
	)
