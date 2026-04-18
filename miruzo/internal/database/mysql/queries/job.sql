-- name: MarkJobStarted :execrows
INSERT INTO jobs(name, started_at) VALUES(?, ?)
ON DUPLICATE KEY UPDATE
	started_at = IF(finished_at IS NOT NULL, VALUES(started_at), started_at),
	finished_at = IF(finished_at IS NOT NULL, NULL, finished_at);

-- name: MarkJobFinished :execrows
UPDATE jobs SET finished_at=? WHERE name=? AND finished_at IS NULL;
