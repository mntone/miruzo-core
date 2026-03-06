-- name: ListImagesEngaged :many
SELECT sqlc.embed(images), stats.score_evaluated
FROM images
JOIN stats ON stats.ingest_id = images.ingest_id
WHERE stats.hall_of_fame_at IS NULL AND stats.score_evaluated >= sqlc.arg(score_threshold)
ORDER BY stats.score_evaluated DESC, images.ingest_id DESC
LIMIT $1;

-- name: ListImagesEngagedAfter :many
SELECT sqlc.embed(images), stats.score_evaluated
FROM images
JOIN stats ON stats.ingest_id = images.ingest_id
WHERE stats.hall_of_fame_at IS NULL AND stats.score_evaluated >= sqlc.arg(score_threshold) AND stats.score_evaluated < $1
ORDER BY stats.score_evaluated DESC, images.ingest_id DESC
LIMIT $2;
