-- name: ApplyHallOfFameGrantedToStats :execrows
UPDATE stats SET hall_of_fame_at=?
WHERE ingest_id=? AND score>=sqlc.arg(hall_of_fame_score_threshold) AND hall_of_fame_at IS NULL;

-- name: ApplyHallOfFameRevokedToStats :execrows
UPDATE stats SET hall_of_fame_at=NULL
WHERE ingest_id=? AND hall_of_fame_at IS NOT NULL;
