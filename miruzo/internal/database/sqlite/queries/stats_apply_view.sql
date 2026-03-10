-- name: ApplyViewToStats :execrows
UPDATE stats
SET
	score=score+sqlc.arg(score_delta),
	last_viewed_at=sqlc.arg(evaluated_at),
	view_count=view_count+1
WHERE stats.ingest_id=?1;

-- name: ApplyViewToStatsWithMilestone :execrows
UPDATE stats
SET
	score=score+sqlc.arg(score_delta),
	last_viewed_at=sqlc.arg(evaluated_at),
	view_count=view_count+1,
	view_milestone_count=view_count+1,
	view_milestone_archived_at=sqlc.arg(evaluated_at)
WHERE stats.ingest_id=?1;
