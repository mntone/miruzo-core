-- name: ApplyHallOfFameGrantedToStats :execrows
UPDATE stats SET hall_of_fame_at=$2
WHERE ingest_id=$1 AND score>=sqlc.arg(hall_of_fame_score_threshold) AND hall_of_fame_at IS NULL;

-- name: ApplyHallOfFameRevokedToStats :execrows
UPDATE stats SET hall_of_fame_at=NULL
WHERE ingest_id=$1 AND hall_of_fame_at IS NOT NULL;
