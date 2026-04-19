package postgres

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

const (
	actionCreateOperationName           = "Create"
	actionCreateDailyDecayOperationName = "CreateDailyDecayIfAbsent"
	actionCreateLoveOperationName       = "CreateLoveIfAbsent"
	actionCreateHallOfFameOperationName = "CreateHallOfFameIfAbsent"
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
) error {
	affectedRows, err := repo.queries.CreateAction(ctx, gen.CreateActionParams{
		IngestID:      ingestID,
		Kind:          int16(kind),
		OccurredAt:    occurredAt,
		PeriodStartAt: periodStartAt,
	})
	if err != nil {
		return dberrors.ToPersist(actionCreateOperationName, err)
	}
	if affectedRows == 1 {
		return nil
	}

	var baseError error
	if affectedRows == 0 {
		baseError = persist.ErrActionAlreadyExists
	} else {
		baseError = persist.ErrInvariantViolation
	}
	return dberrors.WrapKV(
		baseError,
		actionCreateOperationName,
		"affected_rows", affectedRows,
		"ingest_id", ingestID,
		"occurred_at", occurredAt.Format(time.RFC3339Nano),
		"period_start_at", periodStartAt.Format(time.RFC3339Nano),
	)
}

func (repo actionRepository) CreateDailyDecayIfAbsent(
	ctx context.Context,
	ingestID model.IngestIDType,
	occurredAt time.Time,
	periodStartAt time.Time,
) error {
	affectedRows, err := repo.queries.CreateDailyDecayActionIfAbsent(ctx, gen.CreateDailyDecayActionIfAbsentParams{
		IngestID:      ingestID,
		OccurredAt:    occurredAt,
		PeriodStartAt: periodStartAt,
	})
	if err != nil {
		return dberrors.ToPersist(actionCreateDailyDecayOperationName, err)
	}
	if affectedRows == 1 {
		return nil
	}

	var baseError error
	if affectedRows == 0 {
		baseError = persist.ErrActionAlreadyExists
	} else {
		baseError = persist.ErrInvariantViolation
	}
	return dberrors.WrapKV(
		baseError,
		actionCreateDailyDecayOperationName,
		"affected_rows", affectedRows,
		"ingest_id", ingestID,
		"occurred_at", occurredAt.Format(time.RFC3339Nano),
		"period_start_at", periodStartAt.Format(time.RFC3339Nano),
	)
}

func (repo actionRepository) CreateLoveIfAbsent(
	ctx context.Context,
	ingestID model.IngestIDType,
	loveType persist.LoveActionType,
	occurredAt time.Time,
	periodStartAt time.Time,
) error {
	affectedRows, err := repo.queries.CreateLoveActionIfAbsent(ctx, gen.CreateLoveActionIfAbsentParams{
		IngestID:      ingestID,
		Kind:          int16(loveType),
		OccurredAt:    occurredAt,
		PeriodStartAt: periodStartAt,
	})
	if err != nil {
		return dberrors.ToPersist(actionCreateLoveOperationName, err)
	}
	if affectedRows == 1 {
		return nil
	}

	var baseError error
	if affectedRows == 0 {
		baseError = persist.ErrActionAlreadyExists
	} else {
		baseError = persist.ErrInvariantViolation
	}
	return dberrors.WrapKV(
		baseError,
		actionCreateLoveOperationName,
		"affected_rows", affectedRows,
		"ingest_id", ingestID,
		"occurred_at", occurredAt.Format(time.RFC3339Nano),
		"period_start_at", periodStartAt.Format(time.RFC3339Nano),
	)
}

func (repo actionRepository) CreateHallOfFameIfAbsent(
	ctx context.Context,
	ingestID model.IngestIDType,
	hallOfFameType persist.HallOfFameActionType,
	occurredAt time.Time,
	periodStartAt time.Time,
) error {
	affectedRows, err := repo.queries.CreateHallOfFameActionIfAbsent(ctx, gen.CreateHallOfFameActionIfAbsentParams{
		IngestID:      ingestID,
		Kind:          int16(hallOfFameType),
		OccurredAt:    occurredAt,
		PeriodStartAt: periodStartAt,
	})
	if err != nil {
		return dberrors.ToPersist(actionCreateHallOfFameOperationName, err)
	}
	if affectedRows == 1 {
		return nil
	}

	var baseError error
	if affectedRows == 0 {
		baseError = persist.ErrActionAlreadyExists
	} else {
		baseError = persist.ErrInvariantViolation
	}
	return dberrors.WrapKV(
		baseError,
		actionCreateHallOfFameOperationName,
		"affected_rows", affectedRows,
		"ingest_id", ingestID,
		"occurred_at", occurredAt.Format(time.RFC3339Nano),
		"period_start_at", periodStartAt.Format(time.RFC3339Nano),
	)
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
		return false, dberrors.ToPersist("ExistsSince", err)
	}

	return exists, nil
}
