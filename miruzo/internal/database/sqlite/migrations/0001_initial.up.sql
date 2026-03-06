-- Create ingests table
CREATE TABLE ingests(
	id INTEGER PRIMARY KEY AUTOINCREMENT
		CONSTRAINT c_id
			CHECK (id >= 1 AND id <= 9007199254740991),
	process INTEGER
		CONSTRAINT c_process
			NOT NULL
			CHECK (process IN (0, 1))
			DEFAULT 0,
	visibility INTEGER
		CONSTRAINT c_visibility
			NOT NULL
			CHECK (visibility IN (0, 1))
			DEFAULT 0,
	relative_path VARCHAR
		CONSTRAINT c_relative_path
			NOT NULL
			CHECK (length(relative_path) >= 4 AND relative_path NOT LIKE '/%'),
	fingerprint VARCHAR(64) NOT NULL,
	ingested_at DATETIME NOT NULL,
	captured_at DATETIME
		CONSTRAINT c_captured_at
			NOT NULL
			CHECK (captured_at <= ingested_at),
	updated_at DATETIME
		CONSTRAINT c_updated_at
			NOT NULL
			CHECK (updated_at >= ingested_at),
	executions JSON
		CONSTRAINT c_executions
			NOT NULL
			CHECK (json_valid(executions) AND json_type(executions) = 'array')
			DEFAULT '[]',
	UNIQUE (fingerprint)
);

-- Create images table
CREATE TABLE images(
	ingest_id INTEGER
		CONSTRAINT c_images_ingest_id
			PRIMARY KEY
			REFERENCES ingests(id),
	ingested_at DATETIME NOT NULL,
	kind INTEGER
		CONSTRAINT c_kind
			NOT NULL
			CHECK (kind IN (0, 1, 2, 3))
			DEFAULT 0,
	original JSON
		CONSTRAINT c_original
			NOT NULL
			CHECK (json_valid(original) AND json_type(original) = 'object'),
	fallback JSON
		CONSTRAINT c_fallback
			NOT NULL
			CHECK (json_valid(fallback) AND json_type(fallback) IN ('null', 'object')),
	variants JSON
		CONSTRAINT c_variants
			NOT NULL
			CHECK (json_valid(variants) AND json_type(variants) = 'array')
);

-- Create stats table
CREATE TABLE stats(
	ingest_id INTEGER
		CONSTRAINT c_stats_ingest_id
			PRIMARY KEY
			REFERENCES ingests(id),
	score INTEGER
		CONSTRAINT c_score
			NOT NULL
			CHECK (score >= -32768 AND score <= 32767),
	score_evaluated INTEGER
		CONSTRAINT c_score_evaluated
			NOT NULL
			CHECK (score_evaluated >= -32768 AND score_evaluated <= 32767),
	score_evaluated_at DATETIME,
	first_loved_at DATETIME,
	last_loved_at DATETIME,
	hall_of_fame_at DATETIME,
	last_viewed_at DATETIME
		CONSTRAINT c_last_viewed_at
			CHECK (
				(view_count != 0 AND last_viewed_at IS NOT NULL)
				OR
				(view_count == 0 AND last_viewed_at IS NULL)
			),
	view_count INTEGER
		CONSTRAINT c_view_count
			NOT NULL
			CHECK (view_count >= 0)
			DEFAULT 0,
	view_milestone_count INTEGER
		CONSTRAINT c_view_milestone_count
			NOT NULL
			CHECK (view_milestone_count >= 0 AND view_milestone_count <= view_count)
			DEFAULT 0,
	view_milestone_archived_at DATETIME
		CONSTRAINT c_view_milestone_archived_at
			CHECK (
				(view_milestone_count != 0 AND view_milestone_archived_at IS NOT NULL)
				OR
				(view_milestone_count == 0 AND view_milestone_archived_at IS NULL)
			),
	CONSTRAINT c_loved_at_pair
		CHECK (
			(first_loved_at IS NULL AND last_loved_at IS NULL)
			OR
			(first_loved_at IS NOT NULL AND last_loved_at IS NOT NULL)
		),
	CONSTRAINT c_loved_at_order
		CHECK (first_loved_at IS NULL OR first_loved_at <= last_loved_at)
);
