-- name: ListImagesHallOfFame :many
SELECT sqlc.embed(images), stats.hall_of_fame_at
FROM images JOIN stats USING(ingest_id)
WHERE stats.hall_of_fame_at IS NOT NULL
ORDER BY stats.hall_of_fame_at DESC, ingest_id DESC
LIMIT ?;

-- name: ListImagesHallOfFameAfter :many
SELECT sqlc.embed(images), stats.hall_of_fame_at
FROM images JOIN stats USING(ingest_id)
WHERE stats.hall_of_fame_at IS NOT NULL
  AND (
    stats.hall_of_fame_at < sqlc.arg(cursor_at)
    OR
    (stats.hall_of_fame_at = sqlc.arg(cursor_at) AND images.ingest_id < sqlc.arg(cursor_id))
  )
ORDER BY stats.hall_of_fame_at DESC, ingest_id DESC
LIMIT ?;
