-- name: CreateAction :execlastid
INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES(?, ?, ?, ?);

-- name: CreateDailyDecayActionIfAbsent :execrows
INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES(?, 1/*decay*/, ?, ?)
ON CONFLICT(ingest_id, period_start_at) WHERE kind=1 DO NOTHING;

-- name: CreateLoveActionIfAbsent :execrows
INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES(?, ?, ?, ?)
ON CONFLICT(ingest_id, occurred_at) WHERE kind IN(13/*love*/, 14/*love_canceled*/)DO NOTHING;
