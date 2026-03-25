-- name: ListImagesLatest :many
SELECT *
FROM images
ORDER BY ingested_at DESC, ingest_id DESC
LIMIT $1;

-- name: ListImagesLatestAfter :many
SELECT *
FROM images
WHERE (ingested_at, ingest_id) < (sqlc.arg(cursor_at), sqlc.arg(cursor_id)::bigint)
ORDER BY ingested_at DESC, ingest_id DESC
LIMIT sqlc.arg(max_count);
