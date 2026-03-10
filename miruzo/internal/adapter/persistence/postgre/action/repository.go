package action

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
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
	ingestID model.IngestIDType,
	kind model.ActionType,
	occurredAt time.Time,
) (model.ActionIDType, error) {
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
