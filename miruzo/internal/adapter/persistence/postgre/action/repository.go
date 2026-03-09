package action

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model/action"
)

type repository struct {
	queries *gen.Queries
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		queries: gen.New(pool),
	}
}

func (repo *repository) CreateAction(
	ctx context.Context,
	ingestID int64,
	kind action.ActionType,
	occurredAt time.Time,
) (action.ActionID, error) {
	actionID, err := repo.queries.CreateAction(ctx, gen.CreateActionParams{
		IngestID:   ingestID,
		Kind:       int16(kind),
		OccurredAt: shared.PgtypeTimestampFromTime(occurredAt),
	})
	if err != nil {
		return 0, shared.MapPostgreError("CreateAction", err)
	}

	return actionID, nil
}
