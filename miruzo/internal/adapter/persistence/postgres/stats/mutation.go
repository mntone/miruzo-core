package stats

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
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
		EvaluatedAt: &evaluatedAt,
	})
	if err != nil {
		return shared.MapPostgreError("ApplyDecay", err)
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
		HallOfFameAt:             &hallOfFameAt,
		HallOfFameScoreThreshold: hallOfFameScoreThreshold,
	})
	if err != nil {
		return shared.MapPostgreError("ApplyHallOfFameGranted", err)
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
		return shared.MapPostgreError("ApplyHallOfFameRevoked", err)
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
		LovedAt:            &lovedAt,
		PeriodStartAt:      &periodStartAt,
		LoveScoreThreshold: loveScoreThreshold,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.LoveStats{}, persist.ErrConflict
		}

		return model.LoveStats{}, shared.MapPostgreError("ApplyLove", err)
	}

	return model.LoveStats{
		Score:        loveStats.Score,
		FirstLovedAt: mo.PointerToOption(loveStats.FirstLovedAt),
		LastLovedAt:  mo.PointerToOption(loveStats.LastLovedAt),
	}, nil
}

func (repo repository) ApplyLoveCanceled(
	ctx context.Context,
	ingestID model.IngestIDType,
	scoreDelta model.ScoreType,
	loveCanceledAt time.Time,
	periodStartAt time.Time,
	dayStartOffset time.Duration,
) (model.LoveStats, error) {
	loveStats, err := repo.queries.ApplyLoveCanceledToStats(ctx, gen.ApplyLoveCanceledToStatsParams{
		IngestID:       ingestID,
		ScoreDelta:     scoreDelta,
		PeriodStartAt:  &periodStartAt,
		LoveCanceledAt: &loveCanceledAt,
		DayStartOffset: shared.PgtypeIntervalFromDuration(dayStartOffset),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.LoveStats{}, persist.ErrConflict
		}

		return model.LoveStats{}, shared.MapPostgreError("ApplyLoveCanceled", err)
	}

	return model.LoveStats{
		Score:        loveStats.Score,
		FirstLovedAt: mo.PointerToOption(loveStats.FirstLovedAt),
		LastLovedAt:  mo.PointerToOption(loveStats.LastLovedAt),
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
		ViewedAt:   &viewedAt,
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
		ViewedAt:   &viewedAt,
	})
	if err != nil {
		return shared.MapPostgreError("ApplyViewWithMilestone", err)
	}

	if rowCount == 0 {
		return persist.ErrNotFound
	}

	return nil
}
