-- name: ListImagesChronological :many
SELECT sqlc.embed(images), ingests.captured_at
FROM ingests JOIN images ON images.ingest_id = ingests.id
ORDER BY ingests.captured_at DESC, ingests.id DESC
LIMIT $1;

-- name: ListImagesChronologicalAfter :many
SELECT sqlc.embed(images), ingests.captured_at
FROM ingests JOIN images ON images.ingest_id = ingests.id
WHERE (ingests.captured_at, ingests.id) < (sqlc.arg(cursor_at), sqlc.arg(cursor_id)::bigint)
ORDER BY ingests.captured_at DESC, ingests.id DESC
LIMIT sqlc.arg(max_count);
