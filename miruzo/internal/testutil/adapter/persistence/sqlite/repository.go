package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/model/media"
	"github.com/samber/mo"
)

type repository struct {
	queries *gen.Queries
}

func newRepository(db *sql.DB) repository {
	return repository{
		queries: gen.New(db),
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
	original media.Variant,
	fallback mo.Option[media.Variant],
	variants []media.Variant,
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
		Score:          int64(score),
		ScoreEvaluated: int64(scoreEvaluated),
		FirstLovedAt:   shared.NullTimeFromOption(firstLovedAt),
		LastLovedAt:    shared.NullTimeFromOption(lastLovedAt),
		HallOfFameAt:   shared.NullTimeFromOption(hallOfFameAt),
		LastViewedAt:   shared.NullTimeFromOption(lastViewedAt),
		ViewCount:      viewCount,
	})
}
