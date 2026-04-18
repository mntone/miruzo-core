-- name: ListStatsAfterIngestID :many
SELECT ingest_id, score, last_viewed_at
FROM stats
WHERE ingest_id > sqlc.arg(last_ingest_id)
ORDER BY ingest_id
LIMIT ?;
