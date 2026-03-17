package stats

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
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
		LovedAt:       shared.NullTimeFromTime(lovedAt),
		PeriodStartAt: shared.NullTimeFromTime(periodStartAt),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return persist.LoveStats{}, persist.ErrConflict
		}

		return persist.LoveStats{}, shared.MapSQLiteError("ApplyLove", err)
	}

	return persist.LoveStats{
		Score:        model.ScoreType(loveStats.Score),
		FirstLovedAt: shared.OptionTimeFromSql(loveStats.FirstLovedAt),
		LastLovedAt:  shared.OptionTimeFromSql(loveStats.LastLovedAt),
	}, nil
}

func (repo repository) ApplyLoveCanceled(
	ctx context.Context,
	ingestID model.IngestIDType,
	scoreDelta model.ScoreType,
	periodStartAt time.Time,
	dayStartOffset time.Duration,
) (persist.LoveStats, error) {
	loveStats, err := repo.queries.ApplyLoveCanceledToStats(ctx, gen.ApplyLoveCanceledToStatsParams{
		IngestID:       ingestID,
		ScoreDelta:     scoreDelta,
		PeriodStartAt:  shared.NullTimeFromTime(periodStartAt),
		DayStartOffset: int64(dayStartOffset.Seconds()),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return persist.LoveStats{}, persist.ErrConflict
		}

		return persist.LoveStats{}, shared.MapSQLiteError("ApplyLoveCanceled", err)
	}

	return persist.LoveStats{
		Score:        model.ScoreType(loveStats.Score),
		FirstLovedAt: shared.OptionTimeFromSql(loveStats.FirstLovedAt),
		LastLovedAt:  shared.OptionTimeFromSql(loveStats.LastLovedAt),
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
		ViewedAt:   shared.NullTimeFromTime(viewedAt),
	})
	if err != nil {
		return shared.MapSQLiteError("ApplyView", err)
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
		ViewedAt:   shared.NullTimeFromTime(viewedAt),
	})
	if err != nil {
		return shared.MapSQLiteError("ApplyViewWithMilestone", err)
	}

	if rowCount == 0 {
		return persist.ErrNotFound
	}

	return nil
}
