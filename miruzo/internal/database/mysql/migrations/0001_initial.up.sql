-- Create settings table
CREATE TABLE settings(
	`key` CHAR(8) PRIMARY KEY
		CONSTRAINT ck_settings_key
			CHECK (RTRIM(`key`) REGEXP '^[a-z0-9_]{2,8}$'),
	value VARCHAR(255) NOT NULL
);

-- Create ingests table
-- details: docs/database.md
CREATE TABLE ingests(
	id BIGINT AUTO_INCREMENT PRIMARY KEY,
	process TINYINT NOT NULL DEFAULT 0
		CONSTRAINT ck_ingests_process
			CHECK (process IN (0, 1)),
	visibility TINYINT NOT NULL DEFAULT 0
		CONSTRAINT ck_ingests_visibility
			CHECK (visibility IN (0, 1)),
	relative_path VARCHAR(255) NOT NULL
		CONSTRAINT ck_ingests_relative_path
			CHECK (
				char_length(relative_path) >= 5
				AND relative_path NOT LIKE '/%'
				AND relative_path NOT LIKE './%'
				AND relative_path NOT LIKE '..%'
				AND relative_path NOT LIKE '%//%'
				AND relative_path NOT LIKE '%/./%'
				AND relative_path NOT LIKE '%/../%'
			),
	fingerprint VARCHAR(64) NOT NULL,
	ingested_at DATETIME(6) NOT NULL,
	captured_at DATETIME(6) NOT NULL,
	updated_at  DATETIME(6) NOT NULL,
	executions JSON NOT NULL DEFAULT ('[]')
		CONSTRAINT ck_ingests_executions
			CHECK (JSON_TYPE(executions) = 'ARRAY' AND JSON_LENGTH(executions) <= 5),
	UNIQUE uq_ingests_relative_path (relative_path),
	UNIQUE uq_ingests_fingerprint (fingerprint),
	CONSTRAINT ck_ingests_captured_at CHECK (captured_at <= ingested_at),
	CONSTRAINT ck_ingests_updated_at  CHECK (updated_at >= ingested_at)
);

-- Create images table
CREATE TABLE images(
	ingest_id BIGINT PRIMARY KEY,
	ingested_at DATETIME(6) NOT NULL,
	kind TINYINT NOT NULL DEFAULT 0
		CONSTRAINT ck_images_kind
			CHECK (kind IN (0, 1, 2, 3)),
	original JSON NOT NULL
		CONSTRAINT ck_images_original
			CHECK (JSON_TYPE(original) = 'OBJECT'),
	fallback JSON
		CONSTRAINT ck_images_fallback
			CHECK (JSON_TYPE(fallback) = 'OBJECT'),
	variants JSON NOT NULL
		CONSTRAINT ck_images_variants
			CHECK (JSON_TYPE(variants) = 'ARRAY'),
	CONSTRAINT fk_images_ingest
		FOREIGN KEY (ingest_id) REFERENCES ingests(id)
);

-- Create stats table
CREATE TABLE stats(
	ingest_id BIGINT PRIMARY KEY,
	score SMALLINT NOT NULL,
	score_evaluated SMALLINT NOT NULL,
	score_evaluated_at DATETIME(6),
	first_loved_at DATETIME(6),
	last_loved_at DATETIME(6),
	hall_of_fame_at DATETIME(6),
	last_viewed_at DATETIME(6),
	view_count BIGINT NOT NULL DEFAULT 0 CHECK (view_count >= 0),
	view_milestone_count BIGINT NOT NULL DEFAULT 0,
	view_milestone_archived_at DATETIME(6),
	CONSTRAINT fk_stats_ingest
		FOREIGN KEY (ingest_id) REFERENCES ingests(id)
		ON DELETE CASCADE,
	CONSTRAINT ck_stats_loved_at_pair
		CHECK (
			(first_loved_at IS NULL AND last_loved_at IS NULL)
			OR
			(first_loved_at IS NOT NULL AND last_loved_at IS NOT NULL)
		),
	CONSTRAINT ck_stats_loved_at_order
		CHECK (first_loved_at IS NULL OR first_loved_at <= last_loved_at),
	CONSTRAINT ck_stats_last_viewed_at
		CHECK (
			(view_count <> 0 AND last_viewed_at IS NOT NULL)
			OR
			(view_count = 0 AND last_viewed_at IS NULL)
		),
	CONSTRAINT ck_stats_view_milestone_count
		CHECK (view_milestone_count BETWEEN 0 AND view_count),
	CONSTRAINT ck_stats_view_milestone_archived_at
		CHECK (
			(view_milestone_count <> 0 AND view_milestone_archived_at IS NOT NULL)
			OR
			(view_milestone_count = 0 AND view_milestone_archived_at IS NULL)
		)
);

-- Create actions table
CREATE TABLE actions(
	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	ingest_id BIGINT NOT NULL,
	kind TINYINT NOT NULL DEFAULT 0
		CONSTRAINT ck_actions_kind
			CHECK (kind IN (0, 1, 11, 12, 13, 14, 15, 16)),
	occurred_at DATETIME(6) NOT NULL,
	period_start_at DATETIME NOT NULL,
	_decay_period_start_at DATETIME GENERATED ALWAYS AS (
		CASE WHEN kind=1/*decay*/ THEN
			period_start_at
		ELSE NULL
		END
	) VIRTUAL,
	_love_occurred_at DATETIME(6) GENERATED ALWAYS AS (
		CASE WHEN kind IN (13/*love*/, 14/*love_canceled*/) THEN
			occurred_at
		ELSE NULL
		END
	) VIRTUAL,
	_hall_of_fame_occurred_at DATETIME(6) GENERATED ALWAYS AS (
		CASE WHEN kind IN (15/*hall_of_fame_granted*/, 16/*hall_of_fame_revoked*/) THEN
			occurred_at
		ELSE NULL
		END
	) VIRTUAL,
	CONSTRAINT fk_actions_ingest
		FOREIGN KEY (ingest_id) REFERENCES ingests(id)
		ON DELETE CASCADE,
	UNIQUE KEY uq_actions_decay_once_per_period (ingest_id, _decay_period_start_at),
	UNIQUE KEY uq_actions_love_once_per_timestamp (ingest_id, _love_occurred_at),
	UNIQUE KEY uq_actions_hall_of_fame_once_per_timestamp (ingest_id, _hall_of_fame_occurred_at)
);

-- Create jobs table
CREATE TABLE jobs(
	name CHAR(16) PRIMARY KEY
		CONSTRAINT ck_jobs_name
			CHECK (char_length(name) >= 8),
	started_at DATETIME(6) NOT NULL,
	finished_at DATETIME(6)
);

-- Create users table
CREATE TABLE users(
	id SMALLINT PRIMARY KEY AUTO_INCREMENT,
	daily_love_used TINYINT NOT NULL DEFAULT 0
		CONSTRAINT ck_users_daily_love_used
			CHECK (daily_love_used BETWEEN 0 AND 100)
);
INSERT INTO users(id) VALUES(1);
