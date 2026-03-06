-- name: ListImagesFirstLove :many
SELECT sqlc.embed(images), stats.first_loved_at
FROM images
JOIN stats ON stats.ingest_id = images.ingest_id
WHERE stats.first_loved_at IS NOT NULL
ORDER BY stats.first_loved_at DESC, images.ingest_id DESC
LIMIT $1;

-- name: ListImagesFirstLoveAfter :many
SELECT sqlc.embed(images), stats.first_loved_at
FROM images
JOIN stats ON stats.ingest_id = images.ingest_id
WHERE stats.first_loved_at IS NOT NULL AND stats.first_loved_at < $1
ORDER BY stats.first_loved_at DESC, images.ingest_id DESC
LIMIT $2;
