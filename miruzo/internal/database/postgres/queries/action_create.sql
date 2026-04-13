-- name: CreateAction :one
INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES($1, $2, $3, $4)
RETURNING id;

-- name: CreateDailyDecayActionIfAbsent :execrows
INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES($1, 1/*decay*/, $2, $3)
ON CONFLICT(ingest_id, period_start_at) WHERE kind=1 DO NOTHING;
