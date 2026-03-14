-- name: CreateAction :one
INSERT INTO actions(
	ingest_id,
	kind,
	occurred_at
) VALUES($1, $2, $3)
RETURNING id;
