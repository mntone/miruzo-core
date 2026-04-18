-- name: GetLoveStats :one
SELECT score, first_loved_at, last_loved_at FROM stats WHERE ingest_id=?;

-- name: ApplyLoveToStats :execrows
UPDATE stats
SET
	score=score+sqlc.arg(score_delta),
	first_loved_at=COALESCE(first_loved_at, sqlc.arg(loved_at)),
	last_loved_at=sqlc.arg(loved_at)
WHERE ingest_id=?
  AND (last_loved_at IS NULL OR last_loved_at < sqlc.arg(period_start_at))
  AND score<sqlc.arg(love_score_threshold);

-- name: ApplyLoveCanceledToStats :execrows
WITH latest AS (
	SELECT a.occurred_at
	FROM actions a
	JOIN stats s USING(ingest_id)
	WHERE a.ingest_id=sqlc.arg(ingest_id)
	  AND a.kind=13/*love*/
	  AND a.occurred_at >= s.first_loved_at
	  AND a.occurred_at < sqlc.narg(period_start_at)
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
UPDATE stats LEFT JOIN latest ON TRUE
SET
	score=stats.score+sqlc.arg(score_delta),
	first_loved_at=
		CASE
			WHEN latest.occurred_at IS NULL THEN NULL
			WHEN stats.first_loved_at > latest.occurred_at THEN latest.occurred_at
			ELSE stats.first_loved_at
		END,
	last_loved_at=latest.occurred_at
WHERE stats.ingest_id=sqlc.arg(ingest_id)
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
  );
