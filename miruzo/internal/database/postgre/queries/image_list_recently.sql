-- name: ListImagesRecently :many
SELECT sqlc.embed(images), stats.last_viewed_at
FROM images
JOIN stats ON stats.ingest_id = images.ingest_id
WHERE stats.last_viewed_at IS NOT NULL
ORDER BY stats.last_viewed_at DESC, images.ingest_id DESC
LIMIT $1;

-- name: ListImagesRecentlyAfter :many
SELECT sqlc.embed(images), stats.last_viewed_at
FROM images
JOIN stats ON stats.ingest_id = images.ingest_id
WHERE stats.last_viewed_at IS NOT NULL AND stats.last_viewed_at < $1
ORDER BY stats.last_viewed_at DESC, images.ingest_id DESC
LIMIT $2;
