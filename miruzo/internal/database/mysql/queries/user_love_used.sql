-- name: GetDailyLoveUsed :one
SELECT daily_love_used FROM users WHERE id=1;

-- name: IncrementDailyLoveUsed :execrows
UPDATE users
SET daily_love_used=daily_love_used+1
WHERE id=1 AND daily_love_used < sqlc.arg(daily_love_limit);

-- name: DecrementDailyLoveUsed :execrows
UPDATE users SET daily_love_used=daily_love_used-1 WHERE id=1;

-- name: ResetDailyLoveUsed :execrows
UPDATE users SET daily_love_used=0 WHERE id=1;
