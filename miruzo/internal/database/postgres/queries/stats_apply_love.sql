-- name: ApplyLoveToStats :one
UPDATE stats
SET
	score=score+sqlc.arg(score_delta),
	first_loved_at=COALESCE(first_loved_at, sqlc.arg(loved_at)),
	last_loved_at=sqlc.arg(loved_at)
WHERE ingest_id=$1
  AND (last_loved_at IS NULL OR last_loved_at < sqlc.arg(period_start_at))
  AND score<sqlc.arg(love_score_threshold)
RETURNING score, first_loved_at, last_loved_at;

-- name: ApplyLoveCanceledToStats :one
WITH latest AS (
	SELECT a.occurred_at
	FROM actions a
	JOIN stats s USING(ingest_id)
	WHERE a.ingest_id=$1
	  AND a.kind=13/*love*/
	  AND a.occurred_at >= s.first_loved_at
	  AND a.occurred_at < sqlc.arg(period_start_at)
	  AND NOT EXISTS (
	    SELECT 1
	    FROM actions b
	    WHERE b.ingest_id=a.ingest_id
	      AND b.kind=14/*love_canceled*/
	      AND b.occurred_at >= a.occurred_at
	      AND b.period_start_at=a.period_start_at
	  )
	ORDER BY a.occurred_at DESC, a.id DESC
	LIMIT 1
)
UPDATE stats
SET
	score=stats.score+sqlc.arg(score_delta),
	first_loved_at=
		CASE
			WHEN l.occurred_at IS NULL THEN NULL
			WHEN stats.first_loved_at > l.occurred_at THEN l.occurred_at
			ELSE stats.first_loved_at
		END,
	last_loved_at=l.occurred_at
FROM (
	SELECT MAX(occurred_at) AS occurred_at
	FROM latest
) l
WHERE stats.ingest_id=$1
  AND stats.last_loved_at IS NOT NULL
  AND stats.last_loved_at >= sqlc.arg(period_start_at)
  AND stats.last_loved_at < sqlc.arg(love_canceled_at)
  -- Defensive guard: fail safely when stats/action rows are inconsistent.
  AND EXISTS (
    SELECT 1
    FROM actions c
    WHERE c.ingest_id=stats.ingest_id
      AND c.kind=13/*love*/
      AND c.occurred_at=stats.last_loved_at
  )
RETURNING stats.score, stats.first_loved_at, stats.last_loved_at;
