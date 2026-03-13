package stats

import (
	"context"
	"time"

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
