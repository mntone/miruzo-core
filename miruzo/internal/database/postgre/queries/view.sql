-- name: GetImageWithStats :one
SELECT sqlc.embed(images), sqlc.embed(stats)
FROM images JOIN stats USING(ingest_id)
WHERE images.ingest_id=$1;
