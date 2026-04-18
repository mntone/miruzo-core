-- name: ListImagesEngaged :many
SELECT sqlc.embed(images), stats.score_evaluated
FROM images JOIN stats USING(ingest_id)
WHERE stats.hall_of_fame_at IS NULL AND stats.score_evaluated >= sqlc.arg(score_threshold)
ORDER BY stats.score_evaluated DESC, ingest_id DESC
LIMIT ?;

-- name: ListImagesEngagedAfter :many
SELECT sqlc.embed(images), stats.score_evaluated
FROM images JOIN stats USING(ingest_id)
WHERE stats.hall_of_fame_at IS NULL
  AND stats.score_evaluated >= sqlc.arg(score_threshold)
  AND (
    stats.score_evaluated < sqlc.arg(cursor_int)
    OR
    (stats.score_evaluated = sqlc.arg(cursor_int) AND images.ingest_id < sqlc.arg(cursor_id))
  )
ORDER BY stats.score_evaluated DESC, ingest_id DESC
LIMIT ?;
