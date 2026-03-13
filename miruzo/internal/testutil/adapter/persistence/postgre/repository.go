package postgre

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/model/media"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

type repository struct {
	queries *gen.Queries
}

func newRepository(pool *pgxpool.Pool) repository {
	return repository{
		queries: gen.New(pool),
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
		IngestedAt:   shared.PgtypeTimestampFromTime(ingestedAt),
		CapturedAt:   shared.PgtypeTimestampFromTime(capturedAt),
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
	original media.Variant,
	fallback mo.Option[media.Variant],
	variants []media.Variant,
) error {
	err := repo.queries.CreateImage(ctx, gen.CreateImageParams{
		IngestID:   id,
		IngestedAt: shared.PgtypeTimestampFromTime(ingestedAt),
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
		FirstLovedAt:   shared.PgtypeTimestampFromOption(firstLovedAt),
		LastLovedAt:    shared.PgtypeTimestampFromOption(lastLovedAt),
		HallOfFameAt:   shared.PgtypeTimestampFromOption(hallOfFameAt),
		LastViewedAt:   shared.PgtypeTimestampFromOption(lastViewedAt),
		ViewCount:      viewCount,
	})
}

func (repo repository) DeleteUser(ctx context.Context) error {
	rowCount, err := repo.queries.DeleteUser(ctx)
	if err != nil {
		return shared.MapPostgreDeleteError("DeleteUser", err)
	}

	if rowCount == 0 {
		return persist.ErrNotFound
	}

	return nil
}
