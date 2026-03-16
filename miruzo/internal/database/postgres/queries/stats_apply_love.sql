-- name: ApplyLoveToStats :one
UPDATE stats
SET
	score=score+sqlc.arg(score_delta),
	first_loved_at=COALESCE(first_loved_at, sqlc.arg(loved_at)),
	last_loved_at=sqlc.arg(loved_at)
WHERE ingest_id=$1
  AND (last_loved_at IS NULL OR last_loved_at < sqlc.arg(period_start_at))
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
	      AND b.occurred_at <
	            DATE_TRUNC('day', a.occurred_at - sqlc.arg(day_start_offset)::interval)
	            + interval '1 day'
	            + sqlc.arg(day_start_offset)::interval
	  )
	ORDER BY a.occurred_at DESC, a.id DESC
	LIMIT 1
)
UPDATE stats
SET
	score=stats.score+sqlc.arg(score_delta),
	first_loved_at=
		CASE
			WHEN (SELECT occurred_at FROM latest) IS NULL THEN NULL
			WHEN stats.first_loved_at IS NULL
			  OR stats.first_loved_at > (SELECT occurred_at FROM latest)
			THEN (SELECT occurred_at FROM latest)
			ELSE stats.first_loved_at
		END,
	last_loved_at=(SELECT occurred_at FROM latest)
WHERE stats.ingest_id=$1
  AND stats.last_loved_at IS NOT NULL
  AND stats.last_loved_at >= sqlc.arg(period_start_at)
RETURNING stats.score, stats.first_loved_at, stats.last_loved_at;
