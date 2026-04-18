-- name: IncrementDailyLoveUsed :one
UPDATE users
SET daily_love_used=daily_love_used+1
WHERE id=1 AND daily_love_used < sqlc.arg(daily_love_limit)
RETURNING daily_love_used;

-- name: DecrementDailyLoveUsed :one
UPDATE users SET daily_love_used=daily_love_used-1 WHERE id=1
RETURNING daily_love_used;

-- name: ResetDailyLoveUsed :execrows
UPDATE users SET daily_love_used=0 WHERE id=1;
