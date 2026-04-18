-- name: GetSettingsValueByKey :one
SELECT value FROM settings WHERE `key`=?;

-- name: UpdateSettingsValueByKey :exec
INSERT INTO settings(`key`, value) VALUES(?, ?)
ON DUPLICATE KEY UPDATE value=VALUES(value);
