-- name: CreateAction :execlastid
INSERT INTO actions(ingest_id, kind, occurred_at, period_start_at) VALUES(?, ?, ?, ?);
