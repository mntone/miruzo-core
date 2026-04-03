-- name: ExistsActionSince :one
SELECT EXISTS (
	SELECT 1 FROM actions WHERE ingest_id=$1 AND kind=$2 AND occurred_at>=sqlc.arg(since_occurred_at)
);
