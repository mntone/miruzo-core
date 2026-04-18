-- name: GetUser :one
SELECT * FROM users WHERE id=1;

-- name: ExistsUser :one
SELECT EXISTS(SELECT 1 FROM users WHERE id=1);
