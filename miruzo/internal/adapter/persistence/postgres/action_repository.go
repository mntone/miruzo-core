package postgres

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type actionRepository struct {
	queries *gen.Queries
}

func (repo actionRepository) Create(
	ctx context.Context,
	ingestID model.IngestIDType,
	kind model.ActionType,
	occurredAt time.Time,
	periodStartAt time.Time,
) (model.ActionIDType, error) {
	actionID, err := repo.queries.CreateAction(ctx, gen.CreateActionParams{
		IngestID:      ingestID,
		Kind:          int16(kind),
		OccurredAt:    occurredAt,
		PeriodStartAt: periodStartAt,
	})
	if err != nil {
		return 0, shared.MapPostgreError("Create", err)
	}

	return actionID, nil
}

func (repo actionRepository) CreateDailyDecayIfAbsent(
	ctx context.Context,
	ingestID model.IngestIDType,
	occurredAt time.Time,
	periodStartAt time.Time,
) error {
	rowCount, err := repo.queries.CreateDailyDecayActionIfAbsent(ctx, gen.CreateDailyDecayActionIfAbsentParams{
		IngestID:      ingestID,
		OccurredAt:    occurredAt,
		PeriodStartAt: periodStartAt,
	})
	if err != nil {
		return shared.MapPostgreError("CreateDailyDecayIfAbsent", err)
	}

	if rowCount == 0 {
		return persist.ErrConflict
	}

	return nil
}

func (repo actionRepository) CreateLoveIfAbsent(
	ctx context.Context,
	ingestID model.IngestIDType,
	loveType persist.LoveActionType,
	occurredAt time.Time,
	periodStartAt time.Time,
) error {
	rowCount, err := repo.queries.CreateLoveActionIfAbsent(ctx, gen.CreateLoveActionIfAbsentParams{
		IngestID:      ingestID,
		Kind:          int16(loveType),
		OccurredAt:    occurredAt,
		PeriodStartAt: periodStartAt,
	})
	if err != nil {
		return shared.MapPostgreError("CreateLoveIfAbsent", err)
	}

	if rowCount == 0 {
		return persist.ErrConflict
	}

	return nil
}

func (repo actionRepository) ExistsSince(
	ctx context.Context,
	ingestID model.IngestIDType,
	kind model.ActionType,
	sinceOccurredAt time.Time,
) (bool, error) {
	exists, err := repo.queries.ExistsActionSince(ctx, gen.ExistsActionSinceParams{
		IngestID:        ingestID,
		Kind:            int16(kind),
		SinceOccurredAt: sinceOccurredAt,
	})
	if err != nil {
		return false, shared.MapPostgreError("ExistsSince", err)
	}

	return exists, nil
}
