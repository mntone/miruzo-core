package gen

import (
	"context"
	"database/sql"
)

const applyLoveCanceledToStats = `WITH latest AS(SELECT a.occurred_at FROM actions a JOIN stats s USING(ingest_id)WHERE a.ingest_id=?1 AND a.kind=13 AND a.occurred_at>=s.first_loved_at AND a.occurred_at<?3 AND NOT EXISTS(SELECT 1 FROM actions b WHERE b.ingest_id=a.ingest_id AND b.kind=14 AND b.occurred_at>=a.occurred_at AND b.occurred_at<datetime(a.occurred_at,'-'||?4||' seconds','start of day','+1 day','+'||?4||' seconds'))ORDER BY a.occurred_at DESC,a.id DESC LIMIT 1)UPDATE stats SET score=stats.score+?2,first_loved_at=CASE WHEN l.occurred_at IS NULL THEN NULL WHEN stats.first_loved_at IS NULL OR stats.first_loved_at>l.occurred_at THEN l.occurred_at ELSE stats.first_loved_at END,last_loved_at=l.occurred_at FROM(SELECT MAX(occurred_at)AS occurred_at FROM latest)l WHERE stats.ingest_id=?1 AND stats.last_loved_at IS NOT NULL AND stats.last_loved_at>=?3 RETURNING stats.score,stats.first_loved_at,stats.last_loved_at`

type ApplyLoveCanceledToStatsParams struct {
	IngestID       int64
	ScoreDelta     int16
	PeriodStartAt  sql.NullTime
	DayStartOffset int64
}

type ApplyLoveCanceledToStatsRow struct {
	Score        int64
	FirstLovedAt sql.NullTime
	LastLovedAt  sql.NullTime
}

func (q *Queries) ApplyLoveCanceledToStats(ctx context.Context, arg ApplyLoveCanceledToStatsParams) (ApplyLoveCanceledToStatsRow, error) {
	row := q.db.QueryRowContext(ctx, applyLoveCanceledToStats, arg.IngestID, arg.ScoreDelta, arg.PeriodStartAt, arg.DayStartOffset)
	var i ApplyLoveCanceledToStatsRow
	err := row.Scan(&i.Score, &i.FirstLovedAt, &i.LastLovedAt)
	return i, err
}
