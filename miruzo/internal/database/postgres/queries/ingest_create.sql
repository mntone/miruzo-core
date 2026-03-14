-- name: CreateIngest :exec
INSERT INTO ingests(
	id,
	relative_path,
	fingerprint,
	ingested_at,
	captured_at,
	updated_at
) VALUES($1, $2, $3, $4, $5, $4);
