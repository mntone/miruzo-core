-- name: ApplyLoveToStats :one
UPDATE stats
SET
	score=score+sqlc.arg(score_delta),
	first_loved_at=COALESCE(first_loved_at, sqlc.arg(loved_at)),
	last_loved_at=sqlc.arg(loved_at)
WHERE
	ingest_id=$1
	AND
	(last_loved_at IS NULL OR last_loved_at < sqlc.arg(period_start_at))
RETURNING score, first_loved_at, last_loved_at;
