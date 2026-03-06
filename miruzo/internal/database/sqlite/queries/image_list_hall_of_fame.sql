-- name: ListImagesHallOfFame :many
SELECT sqlc.embed(images), stats.hall_of_fame_at
FROM images
JOIN stats ON stats.ingest_id = images.ingest_id
WHERE stats.hall_of_fame_at IS NOT NULL
ORDER BY stats.hall_of_fame_at DESC, images.ingest_id DESC
LIMIT ?;

-- name: ListImagesHallOfFameAfter :many
SELECT sqlc.embed(images), stats.hall_of_fame_at
FROM images
JOIN stats ON stats.ingest_id = images.ingest_id
WHERE stats.hall_of_fame_at IS NOT NULL AND stats.hall_of_fame_at < ?
ORDER BY stats.hall_of_fame_at DESC, images.ingest_id DESC
LIMIT ?;
