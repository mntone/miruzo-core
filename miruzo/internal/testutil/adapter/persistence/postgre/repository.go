package postgre

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model/media"
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
	id int64,
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
	id int64,
	ingestedAt time.Time,
	original media.Variant,
	fallback mo.Option[media.Variant],
	variants []media.Variant,
) error {
	originalBytes, err := json.Marshal(original)
	if err != nil {
		return err
	}

	fallbackBytes, err := json.Marshal(fallback)
	if err != nil {
		return err
	}

	variantsBytes, err := json.Marshal(variants)
	if err != nil {
		return err
	}

	err = repo.queries.CreateImage(ctx, gen.CreateImageParams{
		IngestID:   id,
		IngestedAt: shared.PgtypeTimestampFromTime(ingestedAt),
		Original:   originalBytes,
		Fallback:   fallbackBytes,
		Variants:   variantsBytes,
	})
	if err != nil {
		return shared.MapPostgreError("Create", err)
	}

	return nil
}

func (repo repository) CreateStat(
	ctx context.Context,
	id int64,
	score int16,
	scoreEvaluated int16,
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
