-- name: ListImagesLatest :many
SELECT *
FROM images
ORDER BY ingested_at DESC, ingest_id DESC
LIMIT ?;

-- name: ListImagesLatestAfter :many
SELECT *
FROM images
WHERE ingested_at < sqlc.arg(cursor_at)
   OR (ingested_at = sqlc.arg(cursor_at) AND ingest_id < sqlc.arg(cursor_id))
ORDER BY ingested_at DESC, ingest_id DESC
LIMIT ?;
