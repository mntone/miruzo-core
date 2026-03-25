-- name: ListImagesRecently :many
SELECT sqlc.embed(images), stats.last_viewed_at
FROM images JOIN stats USING(ingest_id)
WHERE stats.last_viewed_at IS NOT NULL
ORDER BY stats.last_viewed_at DESC, images.ingest_id DESC
LIMIT $1;

-- name: ListImagesRecentlyAfter :many
SELECT sqlc.embed(images), stats.last_viewed_at
FROM images JOIN stats USING(ingest_id)
WHERE stats.last_viewed_at IS NOT NULL
  AND (stats.last_viewed_at, images.ingest_id) < (sqlc.arg(cursor_at), sqlc.arg(cursor_id)::bigint)
ORDER BY stats.last_viewed_at DESC, images.ingest_id DESC
LIMIT sqlc.arg(max_count);
