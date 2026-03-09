package persist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model/action"
)

type ActionRepository interface {
	CreateAction(
		ctx context.Context,
		ingestID int64,
		kind action.ActionType,
		occurredAt time.Time,
	) (action.ActionID, error)
}
