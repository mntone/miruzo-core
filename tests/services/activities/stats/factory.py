from datetime import datetime

from sqlmodel import Session

from app.models.records import StatsRecord


def _make_stats_record(
	ingest_id: int,
	*,
	score: int,
	view_count: int,
	last_viewed_at: datetime | None,
	hall_of_fame_at: datetime | None,
) -> StatsRecord:
	return StatsRecord(
		ingest_id=ingest_id,
		score=score,
		view_count=view_count,
		last_viewed_at=last_viewed_at,
		first_loved_at=None,
		hall_of_fame_at=hall_of_fame_at,
	)


def build_stats_record(
	ingest_id: int,
	*,
	score: int = 100,
	view_count: int = 0,
	last_viewed_at: datetime | None = None,
	hall_of_fame_at: datetime | None = None,
) -> StatsRecord:
	return _make_stats_record(
		ingest_id=ingest_id,
		score=score,
		view_count=view_count,
		last_viewed_at=last_viewed_at,
		hall_of_fame_at=hall_of_fame_at,
	)


def add_stats_record(
	session: Session,
	ingest_id: int,
	*,
	score: int = 100,
	view_count: int = 0,
	last_viewed_at: datetime | None = None,
	hall_of_fame_at: datetime | None = None,
) -> StatsRecord:
	stats = _make_stats_record(
		ingest_id=ingest_id,
		score=score,
		view_count=view_count,
		last_viewed_at=last_viewed_at,
		hall_of_fame_at=hall_of_fame_at,
	)
	session.add(stats)

	session.commit()
	session.refresh(stats)
	return stats
