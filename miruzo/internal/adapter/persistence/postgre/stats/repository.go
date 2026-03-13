package stats

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type repository struct {
	queries *gen.Queries
}

func NewRepository(queries *gen.Queries) repository {
	return repository{
		queries: queries,
	}
}

func (repo repository) ApplyLove(
	ctx context.Context,
	ingestID model.IngestIDType,
	scoreDelta model.ScoreType,
	lovedAt time.Time,
	periodStartAt time.Time,
) (persist.LoveStats, error) {
	loveStats, err := repo.queries.ApplyLoveToStats(ctx, gen.ApplyLoveToStatsParams{
		IngestID:      ingestID,
		ScoreDelta:    scoreDelta,
		LovedAt:       shared.PgtypeTimestampFromTime(lovedAt),
		PeriodStartAt: shared.PgtypeTimestampFromTime(periodStartAt),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return persist.LoveStats{}, persist.ErrConflict
		}

		return persist.LoveStats{}, shared.MapPostgreError("ApplyLove", err)
	}

	return persist.LoveStats{
		Score:        loveStats.Score,
		FirstLovedAt: shared.OptionTimeFromPgtype(loveStats.FirstLovedAt),
		LastLovedAt:  shared.OptionTimeFromPgtype(loveStats.LastLovedAt),
	}, nil
}

func (repo repository) ApplyView(
	ctx context.Context,
	ingestID model.IngestIDType,
	scoreDelta model.ScoreType,
	viewedAt time.Time,
) error {
	rowCount, err := repo.queries.ApplyViewToStats(ctx, gen.ApplyViewToStatsParams{
		IngestID:   ingestID,
		ScoreDelta: scoreDelta,
		ViewedAt:   shared.PgtypeTimestampFromTime(viewedAt),
	})
	if err != nil {
		return shared.MapPostgreError("ApplyView", err)
	}

	if rowCount == 0 {
		return persist.ErrNotFound
	}

	return nil
}

func (repo repository) ApplyViewWithMilestone(
	ctx context.Context,
	ingestID model.IngestIDType,
	scoreDelta model.ScoreType,
	viewedAt time.Time,
) error {
	rowCount, err := repo.queries.ApplyViewToStatsWithMilestone(ctx, gen.ApplyViewToStatsWithMilestoneParams{
		IngestID:   ingestID,
		ScoreDelta: scoreDelta,
		ViewedAt:   shared.PgtypeTimestampFromTime(viewedAt),
	})
	if err != nil {
		return shared.MapPostgreError("ApplyViewWithMilestone", err)
	}

	if rowCount == 0 {
		return persist.ErrNotFound
	}

	return nil
}
