-- name: ListImagesChronological :many
SELECT sqlc.embed(images), ingests.captured_at
FROM images
JOIN ingests ON ingests.id = images.ingest_id
ORDER BY ingests.captured_at DESC, images.ingest_id DESC
LIMIT ?;

-- name: ListImagesChronologicalAfter :many
SELECT sqlc.embed(images), ingests.captured_at
FROM images
JOIN ingests ON ingests.id = images.ingest_id
WHERE ingests.captured_at < ?
ORDER BY ingests.captured_at DESC, images.ingest_id DESC
LIMIT ?;
