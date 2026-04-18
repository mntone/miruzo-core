-- name: ListImagesRecently :many
SELECT sqlc.embed(images), stats.last_viewed_at
FROM images JOIN stats USING(ingest_id)
WHERE stats.last_viewed_at IS NOT NULL
ORDER BY stats.last_viewed_at DESC, ingest_id DESC
LIMIT ?;

-- name: ListImagesRecentlyAfter :many
SELECT sqlc.embed(images), stats.last_viewed_at
FROM images JOIN stats USING(ingest_id)
WHERE stats.last_viewed_at IS NOT NULL
  AND (
    stats.last_viewed_at < sqlc.arg(cursor_at)
    OR
    (stats.last_viewed_at = sqlc.arg(cursor_at) AND images.ingest_id < sqlc.arg(cursor_id))
  )
ORDER BY stats.last_viewed_at DESC, ingest_id DESC
LIMIT ?;
