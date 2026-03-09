-- name: CreateAction :one
INSERT INTO actions(
	ingest_id,
	kind,
	occurred_at
) VALUES(?, ?, ?)
RETURNING id;
