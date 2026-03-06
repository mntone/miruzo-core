-- name: CreateImage :exec
INSERT INTO images(
	ingest_id,
	ingested_at,
	original,
	fallback,
	variants
) VALUES(?, ?, ?, ?, ?);
