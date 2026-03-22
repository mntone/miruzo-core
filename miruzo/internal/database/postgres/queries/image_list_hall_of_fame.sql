-- name: ListImagesHallOfFame :many
SELECT sqlc.embed(images), stats.hall_of_fame_at
FROM images JOIN stats USING(ingest_id)
WHERE stats.hall_of_fame_at IS NOT NULL
ORDER BY stats.hall_of_fame_at DESC, images.ingest_id DESC
LIMIT $1;

-- name: ListImagesHallOfFameAfter :many
SELECT sqlc.embed(images), stats.hall_of_fame_at
FROM images JOIN stats USING(ingest_id)
WHERE stats.hall_of_fame_at IS NOT NULL AND stats.hall_of_fame_at < $1
ORDER BY stats.hall_of_fame_at DESC, images.ingest_id DESC
LIMIT $2;
