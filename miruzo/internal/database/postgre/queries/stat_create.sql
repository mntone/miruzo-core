-- name: CreateStat :exec
INSERT INTO stats(
	ingest_id,
	score,
	score_evaluated,
	first_loved_at,
	last_loved_at,
	hall_of_fame_at,
	last_viewed_at,
	view_count
) VALUES($1, $2, $3, $4, $5, $6, $7, $8);
