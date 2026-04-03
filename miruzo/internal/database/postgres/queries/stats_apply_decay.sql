-- name: ApplyDecayToStats :execrows
UPDATE stats
SET
	score=sqlc.arg(score),
	score_evaluated=sqlc.arg(score),
	score_evaluated_at=sqlc.arg(evaluated_at)
WHERE ingest_id=$1;
