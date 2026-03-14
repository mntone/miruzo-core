-- name: ListImagesLatest :many
SELECT *
FROM images
ORDER BY ingested_at DESC, ingest_id DESC
LIMIT $1;

-- name: ListImagesLatestAfter :many
SELECT *
FROM images
WHERE ingested_at < $1
ORDER BY ingested_at DESC, ingest_id DESC
LIMIT $2;
