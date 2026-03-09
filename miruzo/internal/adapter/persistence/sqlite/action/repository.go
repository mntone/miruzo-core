package action

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model/action"
)

type repository struct {
	queries *gen.Queries
}

func NewRepository(queries *gen.Queries) repository {
	return repository{
		queries: queries,
	}
}

func (repo repository) CreateAction(
	ctx context.Context,
	ingestID int64,
	kind action.ActionType,
	occurredAt time.Time,
) (action.ActionID, error) {
	actionID, err := repo.queries.CreateAction(ctx, gen.CreateActionParams{
		IngestID:   ingestID,
		Kind:       int64(kind),
		OccurredAt: occurredAt,
	})
	if err != nil {
		return 0, shared.MapSQLiteError("CreateAction", err)
	}

	return actionID, nil
}
