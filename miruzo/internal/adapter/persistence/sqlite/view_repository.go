package sqlite

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/image"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/stats"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type viewRepository struct {
	queries *gen.Queries
}

func NewViewRepository(queries *gen.Queries) viewRepository {
	return viewRepository{
		queries: queries,
	}
}

func (repo viewRepository) GetImageWithStatsForUpdate(
	ctx context.Context,
	ingestID model.IngestIDType,
) (persist.ImageWithStats, error) {
	row, err := repo.queries.GetImageWithStats(ctx, ingestID)
	if err != nil {
		return persist.ImageWithStats{}, shared.MapSQLiteError("GetImageWithStatsForUpdate", err)
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
