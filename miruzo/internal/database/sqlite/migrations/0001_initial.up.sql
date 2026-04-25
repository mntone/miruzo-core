-- Create settings table
CREATE TABLE settings(
	key TEXT PRIMARY KEY
		CONSTRAINT ck_settings_key
			CHECK (length(key) BETWEEN 2 AND 8 AND key NOT GLOB '*[^a-z0-9_]*'),
	value TEXT NOT NULL
		CONSTRAINT ck_settings_value
			NOT NULL
			CHECK (length(value) <= 255)
) STRICT;

-- Create ingests table
-- details: docs/database.md
CREATE TABLE ingests(
	id INTEGER PRIMARY KEY AUTOINCREMENT
		CONSTRAINT ck_ingests_id
			CHECK (id BETWEEN 1 AND 9007199254740991),
	process INTEGER
		CONSTRAINT ck_ingests_process
			NOT NULL
			CHECK (process IN (0, 1))
			DEFAULT 0,
	visibility INTEGER
		CONSTRAINT ck_ingests_visibility
			NOT NULL
			CHECK (visibility IN (0, 1))
			DEFAULT 0,
	relative_path TEXT
		CONSTRAINT ck_ingests_relative_path
			NOT NULL
			CHECK (
				length(relative_path) BETWEEN 5 AND 255
				AND relative_path NOT LIKE '/%'
				AND relative_path NOT LIKE './%'
				AND relative_path NOT LIKE '..%'
				AND relative_path NOT LIKE '%//%'
				AND relative_path NOT LIKE '%/./%'
				AND relative_path NOT LIKE '%/../%'
			)
			UNIQUE,
	fingerprint TEXT
		CONSTRAINT uq_ingests_fingerprint
			NOT NULL
			UNIQUE,
	ingested_at DATETIME NOT NULL,
	captured_at DATETIME NOT NULL CHECK (captured_at <= ingested_at),
	updated_at  DATETIME NOT NULL CHECK (updated_at >= ingested_at),
	executions JSON
		CONSTRAINT ck_ingests_executions
			NOT NULL
			CHECK (json_type(executions) IS 'array' AND json_array_length(executions) <= 5)
			DEFAULT '[]'
);

-- Create images table
CREATE TABLE images(
	ingest_id INTEGER
		CONSTRAINT ck_images_ingest_id
			PRIMARY KEY
			REFERENCES ingests(id),
	ingested_at DATETIME NOT NULL,
	kind INTEGER
		CONSTRAINT ck_images_kind
			NOT NULL
			CHECK (kind IN (0, 1, 2, 3))
			DEFAULT 0,
	original JSON
		CONSTRAINT ck_images_original
			NOT NULL
			CHECK (json_type(original) IS 'object'),
	fallback JSON
		CONSTRAINT ck_images_fallback
			CHECK (fallback IS NULL OR json_type(fallback) IS 'object'),
	variants JSON
		CONSTRAINT ck_images_variants
			NOT NULL
			CHECK (json_type(variants) IS 'array')
);

-- Create stats table
CREATE TABLE stats(
	ingest_id INTEGER
		CONSTRAINT ck_stats_ingest_id
			PRIMARY KEY
			REFERENCES ingests(id)
			ON DELETE CASCADE,
	score INTEGER
		CONSTRAINT ck_stats_score
			NOT NULL
			CHECK (score BETWEEN -32768 AND 32767),
	score_evaluated INTEGER
		CONSTRAINT ck_stats_score_evaluated
			NOT NULL
			CHECK (score_evaluated BETWEEN -32768 AND 32767),
	score_evaluated_at DATETIME,
	first_loved_at DATETIME,
	last_loved_at DATETIME,
	hall_of_fame_at DATETIME,
	last_viewed_at DATETIME
		CONSTRAINT ck_stats_last_viewed_at
			CHECK (
				(view_count != 0 AND last_viewed_at IS NOT NULL)
				OR
				(view_count == 0 AND last_viewed_at IS NULL)
			),
	view_count INTEGER
		CONSTRAINT ck_stats_view_count
			NOT NULL
			CHECK (view_count >= 0)
			DEFAULT 0,
	view_milestone_count INTEGER
		CONSTRAINT ck_stats_view_milestone_count
			NOT NULL
			CHECK (view_milestone_count BETWEEN 0 AND view_count)
			DEFAULT 0,
	view_milestone_archived_at DATETIME
		CONSTRAINT ck_stats_view_milestone_archived_at
			CHECK (
				(view_milestone_count != 0 AND view_milestone_archived_at IS NOT NULL)
				OR
				(view_milestone_count == 0 AND view_milestone_archived_at IS NULL)
			),
	CONSTRAINT ck_stats_loved_at_pair
		CHECK (
			(first_loved_at IS NULL AND last_loved_at IS NULL)
			OR
			(first_loved_at IS NOT NULL AND last_loved_at IS NOT NULL)
		),
	CONSTRAINT ck_stats_loved_at_order
		CHECK (first_loved_at IS NULL OR first_loved_at <= last_loved_at)
);

-- Create actions table
CREATE TABLE actions(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	ingest_id INTEGER
		CONSTRAINT ck_actions_ingest_id
			NOT NULL
			REFERENCES ingests(id)
			ON DELETE CASCADE,
	kind INTEGER
		CONSTRAINT ck_actions_kind
			NOT NULL
			CHECK (kind IN (0, 1, 11, 12, 13, 14, 15, 16))
			DEFAULT 0,
	occurred_at TEXT NOT NULL,
	period_start_at TEXT NOT NULL
) STRICT;
CREATE UNIQUE INDEX uq_actions_decay_once_per_period
ON actions(ingest_id, period_start_at)
WHERE kind=1/*decay*/;
CREATE UNIQUE INDEX uq_actions_love_once_per_timestamp
ON actions(ingest_id, occurred_at)
WHERE kind IN(13/*love*/, 14/*love_canceled*/);
CREATE UNIQUE INDEX uq_actions_hall_of_fame_once_per_timestamp
ON actions(ingest_id, occurred_at)
WHERE kind IN(15/*hall_of_fame_granted*/, 16/*hall_of_fame_revoked*/);

-- Create jobs table
CREATE TABLE jobs(
	name TEXT PRIMARY KEY
		CONSTRAINT ck_jobs_name
			CHECK (length(name) BETWEEN 8 AND 16),
	started_at TEXT NOT NULL,
	finished_at TEXT
) STRICT;

-- Create users table
CREATE TABLE users(
	id INTEGER PRIMARY KEY AUTOINCREMENT
		CONSTRAINT ck_users_id
			CHECK (id BETWEEN 1 AND 32767),
	daily_love_used INTEGER
		CONSTRAINT ck_users_daily_love_used
			NOT NULL
			CHECK (daily_love_used BETWEEN 0 AND 100)
			DEFAULT 0
) STRICT;
INSERT INTO users(id) VALUES(1);
