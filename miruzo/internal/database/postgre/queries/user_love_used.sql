-- name: SetDailyLoveUsed :execrows
UPDATE users SET daily_love_used=$1 WHERE id=1;
