-- name: CreateAction :execrows
INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES(?, ?, ?, ?);

-- name: CreateDailyDecayActionIfAbsent :execrows
INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES(?, 1/*decay*/, ?, ?)
ON DUPLICATE KEY UPDATE ingest_id=ingest_id;

-- name: CreateLoveActionIfAbsent :execrows
INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES(?, ?, ?, ?)
ON DUPLICATE KEY UPDATE ingest_id=ingest_id;

-- name: CreateHallOfFameActionIfAbsent :execrows
INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES(?, ?, ?, ?)
ON DUPLICATE KEY UPDATE ingest_id=ingest_id;
