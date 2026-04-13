-- name: CreateAction :one
INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES($1, $2, $3, $4)
RETURNING id;
