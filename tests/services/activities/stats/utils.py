from datetime import datetime, timezone

from app.models.records import StatsRecord


def _make_stats_record(ingest_id: int) -> StatsRecord:
	return StatsRecord(
		ingest_id=ingest_id,
		hall_of_fame_at=None,
		score=5,
		view_count=1,
		last_viewed_at=datetime.now(timezone.utc),
	)


def build_stats_record(ingest_id: int) -> StatsRecord:
	return _make_stats_record(
		ingest_id=ingest_id,
	)
