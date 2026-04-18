package stats

import (
	"context"
	"database/sql"
	"errors"
	"time"

	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func (repo repository) ApplyDecay(
	ctx context.Context,
	ingestID model.IngestIDType,
	score model.ScoreType,
	evaluatedAt time.Time,
) error {
	rowCount, err := repo.queries.ApplyDecayToStats(ctx, gen.ApplyDecayToStatsParams{
		IngestID:    ingestID,
		Score:       score,
		EvaluatedAt: persistshared.NullTimeFromTime(evaluatedAt),
	})
	if err != nil {
		return dberrors.ToPersist("ApplyDecay", err)
	}

	if rowCount == 0 {
		return persist.ErrConflict
	}

	return nil
}

func (repo repository) ApplyHallOfFameGranted(
	ctx context.Context,
	ingestID model.IngestIDType,
	hallOfFameAt time.Time,
	hallOfFameScoreThreshold model.ScoreType,
) error {
	rowCount, err := repo.queries.ApplyHallOfFameGrantedToStats(ctx, gen.ApplyHallOfFameGrantedToStatsParams{
		IngestID:                 ingestID,
		HallOfFameAt:             persistshared.NullTimeFromTime(hallOfFameAt),
		HallOfFameScoreThreshold: hallOfFameScoreThreshold,
	})
	if err != nil {
		return dberrors.ToPersist("ApplyHallOfFameGranted", err)
	}

	if rowCount == 0 {
		return persist.ErrConflict
	}

	return nil
}

func (repo repository) ApplyHallOfFameRevoked(
	ctx context.Context,
	ingestID model.IngestIDType,
) error {
	rowCount, err := repo.queries.ApplyHallOfFameRevokedToStats(ctx, ingestID)
	if err != nil {
		return dberrors.ToPersist("ApplyHallOfFameRevoked", err)
	}

	if rowCount == 0 {
		return persist.ErrConflict
	}

	return nil
}

func (repo repository) ApplyLove(
	ctx context.Context,
	ingestID model.IngestIDType,
	scoreDelta model.ScoreType,
	lovedAt time.Time,
	loveScoreThreshold model.ScoreType,
	periodStartAt time.Time,
) (model.LoveStats, error) {
	loveStats, err := repo.queries.ApplyLoveToStats(ctx, gen.ApplyLoveToStatsParams{
		IngestID:           ingestID,
		ScoreDelta:         scoreDelta,
		LovedAt:            persistshared.NullTimeFromTime(lovedAt),
		PeriodStartAt:      persistshared.NullTimeFromTime(periodStartAt),
		LoveScoreThreshold: loveScoreThreshold,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.LoveStats{}, persist.ErrConflict
		}

		return model.LoveStats{}, dberrors.ToPersist("ApplyLove", err)
	}

	return model.LoveStats{
		Score:        loveStats.Score,
		FirstLovedAt: persistshared.OptionTimeFromSql(loveStats.FirstLovedAt),
		LastLovedAt:  persistshared.OptionTimeFromSql(loveStats.LastLovedAt),
	}, nil
}

func (repo repository) ApplyLoveCanceled(
	ctx context.Context,
	ingestID model.IngestIDType,
	scoreDelta model.ScoreType,
	loveCanceledAt time.Time,
	periodStartAt time.Time,
) (model.LoveStats, error) {
	loveStats, err := repo.queries.ApplyLoveCanceledToStats(ctx, gen.ApplyLoveCanceledToStatsParams{
		IngestID:       ingestID,
		ScoreDelta:     scoreDelta,
		PeriodStartAt:  persistshared.NullTimeFromTime(periodStartAt),
		LoveCanceledAt: persistshared.NullTimeFromTime(loveCanceledAt),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.LoveStats{}, persist.ErrConflict
		}

		return model.LoveStats{}, dberrors.ToPersist("ApplyLoveCanceled", err)
	}

	return model.LoveStats{
		Score:        loveStats.Score,
		FirstLovedAt: persistshared.OptionTimeFromSql(loveStats.FirstLovedAt),
		LastLovedAt:  persistshared.OptionTimeFromSql(loveStats.LastLovedAt),
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
		ViewedAt:   persistshared.NullTimeFromTime(viewedAt),
	})
	if err != nil {
		return dberrors.ToPersist("ApplyView", err)
	}

	if rowCount == 0 {
		return persist.ErrNoRows
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
		ViewedAt:   persistshared.NullTimeFromTime(viewedAt),
	})
	if err != nil {
		return dberrors.ToPersist("ApplyViewWithMilestone", err)
	}

	if rowCount == 0 {
		return persist.ErrNoRows
	}

	return nil
}
