-- name: ExistsActionSince :one
SELECT EXISTS (
	SELECT 1 FROM actions WHERE ingest_id=? AND kind=? AND occurred_at>=sqlc.arg(since_occurred_at)
);
