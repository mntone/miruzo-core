package postgre

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/image"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/stats"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type repository struct {
	queries *gen.Queries
}

func NewViewRepository(queries *gen.Queries) repository {
	return repository{
		queries: queries,
	}
}

func (repo repository) GetImageWithStats(
	ctx context.Context,
	ingestID model.IngestIDType,
) (persist.ImageWithStats, error) {
	row, err := repo.queries.GetImageWithStats(ctx, ingestID)
	if err != nil {
		return persist.ImageWithStats{}, shared.MapPostgreError("GetImageWithStats", err)
	}

	imageResult, err := image.MapImage(row.Image)
	if err != nil {
		return persist.ImageWithStats{}, err
	}

	return persist.ImageWithStats{
		Image: imageResult,
		Stats: stats.MapStats(row.Stat),
	}, nil
}
