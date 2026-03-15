-- name: GetSettingsValueByKey :one
SELECT value FROM settings WHERE key=?;

-- name: UpdateSettingsValueByKey :exec
INSERT INTO settings(key, value) VALUES(?, ?)
ON CONFLICT(key) DO UPDATE SET value=EXCLUDED.value;
