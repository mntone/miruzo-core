package action

import (
	"context"
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

func (repo repository) Create(
	ctx context.Context,
	ingestID model.IngestIDType,
	kind model.ActionType,
	occurredAt time.Time,
	periodStartAt time.Time,
) (model.ActionIDType, error) {
	actionID, err := repo.queries.CreateAction(ctx, gen.CreateActionParams{
		IngestID:      ingestID,
		Kind:          int64(kind),
		OccurredAt:    occurredAt,
		PeriodStartAt: periodStartAt,
	})
	if err != nil {
		return 0, shared.MapSQLiteError("Create", err)
	}

	return actionID, nil
}

func (repo repository) CreateDailyDecayIfAbsent(
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
		return shared.MapSQLiteError("CreateDailyDecayIfAbsent", err)
	}

	if rowCount == 0 {
		return persist.ErrConflict
	}

	return nil
}

func (repo repository) ExistsSince(
	ctx context.Context,
	ingestID model.IngestIDType,
	kind model.ActionType,
	sinceOccurredAt time.Time,
) (bool, error) {
	exists, err := repo.queries.ExistsActionSince(ctx, gen.ExistsActionSinceParams{
		IngestID:        ingestID,
		Kind:            int64(kind),
		SinceOccurredAt: sinceOccurredAt,
	})
	if err != nil {
		return false, shared.MapSQLiteError("ExistsSince", err)
	}

	return exists != 0, nil
}
