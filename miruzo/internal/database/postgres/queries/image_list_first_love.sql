-- name: ListImagesFirstLove :many
SELECT sqlc.embed(images), stats.first_loved_at
FROM images JOIN stats USING(ingest_id)
WHERE stats.first_loved_at IS NOT NULL
ORDER BY stats.first_loved_at DESC, images.ingest_id DESC
LIMIT $1;

-- name: ListImagesFirstLoveAfter :many
SELECT sqlc.embed(images), stats.first_loved_at
FROM images JOIN stats USING(ingest_id)
WHERE stats.first_loved_at IS NOT NULL
  AND (stats.first_loved_at, images.ingest_id) < (sqlc.arg(cursor_at), sqlc.arg(cursor_id)::bigint)
ORDER BY stats.first_loved_at DESC, images.ingest_id DESC
LIMIT sqlc.arg(max_count);
