package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

type repository struct {
	pool    *pgxpool.Pool
	queries *gen.Queries
}

func newRepository(pool *pgxpool.Pool, queries *gen.Queries) repository {
	return repository{
		pool:    pool,
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
	})
	if err != nil {
		return shared.MapPostgreError("Create", err)
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
	err := repo.queries.CreateImage(ctx, gen.CreateImageParams{
		IngestID:   id,
		IngestedAt: ingestedAt,
		Original:   original,
		Fallback:   fallback.ToPointer(),
		Variants:   variants,
	})
	if err != nil {
		return shared.MapPostgreError("Create", err)
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
		FirstLovedAt:   firstLovedAt.ToPointer(),
		LastLovedAt:    lastLovedAt.ToPointer(),
		HallOfFameAt:   hallOfFameAt.ToPointer(),
		LastViewedAt:   lastViewedAt.ToPointer(),
		ViewCount:      viewCount,
	})
}

func (repo repository) ExecuteStatement(ctx context.Context, stmt string, delete bool) error {
	_, err := repo.pool.Exec(ctx, stmt)
	if err != nil {
		if delete {
			return shared.MapPostgreDeleteError("ExecuteStatement", err)
		}

		return shared.MapPostgreError("ExecuteStatement", err)
	}

	return nil
}

func (repo repository) ExecuteStatementAndReturnRowCount(ctx context.Context, stmt string, delete bool) (int64, error) {
	result, err := repo.pool.Exec(ctx, stmt)
	if err != nil {
		if delete {
			return 0, shared.MapPostgreDeleteError("ExecuteStatementAndReturnRowCount", err)
		}

		return 0, shared.MapPostgreError("ExecuteStatementAndReturnRowCount", err)
	}

	return result.RowsAffected(), nil
}

func (repo repository) TruncateActions(ctx context.Context) error {
	_, err := repo.pool.Exec(ctx, "TRUNCATE TABLE actions")
	if err != nil {
		return shared.MapPostgreDeleteError("TruncateActions", err)
	}

	return nil
}

func (repo repository) TruncateJobs(ctx context.Context) error {
	_, err := repo.pool.Exec(ctx, "TRUNCATE TABLE jobs")
	if err != nil {
		return shared.MapPostgreDeleteError("TruncateJobs", err)
	}

	return nil
}

func (repo repository) TruncateStats(ctx context.Context) error {
	_, err := repo.pool.Exec(ctx, "TRUNCATE TABLE stats")
	if err != nil {
		return shared.MapPostgreDeleteError("TruncateStats", err)
	}

	return nil
}
