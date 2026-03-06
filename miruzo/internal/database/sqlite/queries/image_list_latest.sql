-- name: ListImagesLatest :many
SELECT *
FROM images
ORDER BY ingested_at DESC, ingest_id DESC
LIMIT ?;

-- name: ListImagesLatestAfter :many
SELECT *
FROM images
WHERE ingested_at < ?
ORDER BY ingested_at DESC, ingest_id DESC
LIMIT ?;
