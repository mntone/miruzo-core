package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

type repository struct {
	db      *sql.DB
	queries *gen.Queries
}

func newRepository(db *sql.DB, queries *gen.Queries) repository {
	return repository{
		db:      db,
		queries: queries,
	}
}

func (repo repository) CreateIngest(
	ctx context.Context,
	id model.IngestIDType,
	relativePath string,
	fingerprint string,
	ingestedAt time.Time,
	capturedAt time.Time,
) error {
	err := repo.queries.CreateIngest(ctx, gen.CreateIngestParams{
		ID:           id,
		RelativePath: relativePath,
		Fingerprint:  fingerprint,
		IngestedAt:   ingestedAt,
		CapturedAt:   capturedAt,
		UpdatedAt:    ingestedAt,
	})
	if err != nil {
		return shared.MapSQLiteError("Create", err)
	}

	return nil
}

func (repo repository) CreateImage(
	ctx context.Context,
	id model.IngestIDType,
	ingestedAt time.Time,
	original persist.Variant,
	fallback mo.Option[persist.Variant],
	variants []persist.Variant,
) error {
	originalBytes, err := json.Marshal(original)
	if err != nil {
		return err
	}

	var fallbackBytes *[]byte
	if fallbackValue, present := fallback.Get(); present {
		bytes, err := json.Marshal(fallbackValue)
		if err != nil {
			return err
		}

		fallbackBytes = &bytes
	}

	variantsBytes, err := json.Marshal(variants)
	if err != nil {
		return err
	}

	err = repo.queries.CreateImage(ctx, gen.CreateImageParams{
		IngestID:   id,
		IngestedAt: ingestedAt,
		Original:   originalBytes,
		Fallback:   fallbackBytes,
		Variants:   variantsBytes,
	})
	if err != nil {
		return shared.MapSQLiteError("Create", err)
	}

	return nil
}

func (repo repository) CreateStat(
	ctx context.Context,
	id model.IngestIDType,
	score model.ScoreType,
	scoreEvaluated model.ScoreType,
	lastViewedAt mo.Option[time.Time],
	firstLovedAt mo.Option[time.Time],
	lastLovedAt mo.Option[time.Time],
	hallOfFameAt mo.Option[time.Time],
	viewCount int64,
) error {
	return repo.queries.CreateStat(ctx, gen.CreateStatParams{
		IngestID:       id,
		Score:          score,
		ScoreEvaluated: scoreEvaluated,
		FirstLovedAt:   shared.NullTimeFromOption(firstLovedAt),
		LastLovedAt:    shared.NullTimeFromOption(lastLovedAt),
		HallOfFameAt:   shared.NullTimeFromOption(hallOfFameAt),
		LastViewedAt:   shared.NullTimeFromOption(lastViewedAt),
		ViewCount:      viewCount,
	})
}

func (repo repository) ExecuteStatement(ctx context.Context, stmt string, delete bool) error {
	_, err := repo.db.ExecContext(ctx, stmt)
	if err != nil {
		if delete {
			return shared.MapSQLiteDeleteError("ExecuteStatement", err)
		}

		return shared.MapSQLiteError("ExecuteStatement", err)
	}

	return nil
}

func (repo repository) ExecuteStatementAndReturnRowCount(ctx context.Context, stmt string, delete bool) (int64, error) {
	result, err := repo.db.ExecContext(ctx, stmt)
	if err != nil {
		if delete {
			return 0, shared.MapSQLiteDeleteError("ExecuteStatementAndReturnRowCount", err)
		}

		return 0, shared.MapSQLiteError("ExecuteStatementAndReturnRowCount", err)
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		if delete {
			return 0, shared.MapSQLiteDeleteError("ExecuteStatementAndReturnRowCount", err)
		}

		return 0, shared.MapSQLiteError("ExecuteStatementAndReturnRowCount", err)
	}

	return rowCount, nil
}

func (repo repository) TruncateActions(ctx context.Context) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM actions")
	if err != nil {
		return shared.MapSQLiteDeleteError("TruncateActions", err)
	}

	return nil
}

func (repo repository) TruncateJobs(ctx context.Context) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM jobs")
	if err != nil {
		return shared.MapSQLiteDeleteError("TruncateJobs", err)
	}

	return nil
}

func (repo repository) TruncateStats(ctx context.Context) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM stats")
	if err != nil {
		return shared.MapSQLiteDeleteError("TruncateStats", err)
	}

	return nil
}
