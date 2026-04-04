-- name: MarkJobStarted :execrows
INSERT INTO jobs(name, started_at) VALUES($1, $2) ON CONFLICT(name)
DO UPDATE SET started_at=excluded.started_at, finished_at=NULL
WHERE jobs.name=excluded.name AND jobs.finished_at IS NOT NULL;

-- name: MarkJobFinished :execrows
UPDATE jobs SET finished_at=$2 WHERE name=$1 AND finished_at IS NULL;
