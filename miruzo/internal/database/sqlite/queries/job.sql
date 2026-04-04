-- name: MarkJobStarted :execrows
INSERT INTO jobs(name, started_at) VALUES(?, ?) ON CONFLICT(name)
DO UPDATE SET started_at=EXCLUDED.started_at, finished_at=NULL
WHERE name=EXCLUDED.name AND finished_at IS NOT NULL;

-- name: MarkJobFinished :execrows
UPDATE jobs SET finished_at=? WHERE name=? AND finished_at IS NULL;
