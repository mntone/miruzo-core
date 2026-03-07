-- name: GetSingletonUser :one
SELECT * FROM users
WHERE id = 1
LIMIT 1;
