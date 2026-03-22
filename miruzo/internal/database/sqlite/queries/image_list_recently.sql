-- name: ListImagesRecently :many
SELECT sqlc.embed(images), stats.last_viewed_at
FROM images JOIN stats USING(ingest_id)
WHERE stats.last_viewed_at IS NOT NULL
ORDER BY stats.last_viewed_at DESC, ingest_id DESC
LIMIT ?;

-- name: ListImagesRecentlyAfter :many
SELECT sqlc.embed(images), stats.last_viewed_at
FROM images JOIN stats USING(ingest_id)
WHERE stats.last_viewed_at IS NOT NULL AND stats.last_viewed_at < ?
ORDER BY stats.last_viewed_at DESC, ingest_id DESC
LIMIT ?;
