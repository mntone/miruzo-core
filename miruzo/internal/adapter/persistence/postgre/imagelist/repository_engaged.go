package imagelist

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapEngagedRows(rows []gen.ListImagesEngagedRow) ([]persist.ImageWithCursor[int16], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesEngagedRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesEngagedRow) int16 {
			return row.ScoreEvaluated
		},
	)
}

func mapEngagedAfterRows(rows []gen.ListImagesEngagedAfterRow) ([]persist.ImageWithCursor[int16], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesEngagedAfterRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesEngagedAfterRow) int16 {
			return row.ScoreEvaluated
		},
	)
}

func (repo repository) ListEngaged(
	ctx context.Context,
	spec persist.EngagedImageListSpec,
) ([]persist.ImageWithCursor[int16], error) {
	cursor, present := spec.Cursor.Get()
	if !present {
		rows, err := repo.queries.ListImagesEngaged(
			ctx,
			gen.ListImagesEngagedParams{
				Limit:          int32(spec.Limit),
				ScoreThreshold: spec.ScoreThreshold,
			},
		)
		if err != nil {
			return nil, shared.MapPostgreError("ListEngaged", err)
		}

		return mapEngagedRows(rows)
	}

	rows, err := repo.queries.ListImagesEngagedAfter(
		ctx,
		gen.ListImagesEngagedAfterParams{
			ScoreEvaluated: cursor,
			Limit:          int32(spec.Limit),
			ScoreThreshold: spec.ScoreThreshold,
		},
	)
	if err != nil {
		return nil, shared.MapPostgreError("ListEngaged", err)
	}

	return mapEngagedAfterRows(rows)
}
