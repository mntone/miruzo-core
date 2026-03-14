-- Create index for imagelist
CREATE INDEX ix_images_latest
ON images (ingested_at DESC, ingest_id DESC);

CREATE INDEX ix_ingests_chronological
ON ingests (captured_at DESC, id DESC);

CREATE INDEX ix_stats_recently
ON stats (last_viewed_at DESC, ingest_id DESC)
WHERE last_viewed_at IS NOT NULL;

CREATE INDEX ix_stats_first_love
ON stats (first_loved_at DESC, ingest_id DESC)
WHERE first_loved_at IS NOT NULL;

CREATE INDEX ix_stats_hall_of_fame
ON stats (hall_of_fame_at DESC, ingest_id DESC)
WHERE hall_of_fame_at IS NOT NULL;

CREATE INDEX ix_stats_engaged
ON stats (score_evaluated DESC, ingest_id DESC)
WHERE hall_of_fame_at IS NULL;
