package persist

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type ViewRepository interface {
	GetImageWithStatsForUpdate(
		ctx context.Context,
		ingestID model.IngestIDType,
	) (ImageWithStats, error)
}
