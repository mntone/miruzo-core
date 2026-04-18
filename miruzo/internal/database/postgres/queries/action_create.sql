-- name: CreateAction :execrows
INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES($1, $2, $3, $4);

-- name: CreateDailyDecayActionIfAbsent :execrows
INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES($1, 1/*decay*/, $2, $3)
ON CONFLICT(ingest_id, period_start_at) WHERE kind=1 DO NOTHING;

-- name: CreateLoveActionIfAbsent :execrows
INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES($1, $2, $3, $4)
ON CONFLICT(ingest_id, occurred_at) WHERE kind IN(13/*love*/, 14/*love_canceled*/)DO NOTHING;

-- name: CreateHallOfFameActionIfAbsent :execrows
INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES($1, $2, $3, $4)
ON CONFLICT(ingest_id, occurred_at) WHERE kind IN(15/*hall_of_fame_granted*/, 16/*hall_of_fame_revoked*/)DO NOTHING;
