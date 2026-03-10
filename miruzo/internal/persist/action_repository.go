package persist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type ActionRepository interface {
	CreateAction(
		ctx context.Context,
		ingestID model.IngestIDType,
		kind model.ActionType,
		occurredAt time.Time,
	) (model.ActionIDType, error)
}
