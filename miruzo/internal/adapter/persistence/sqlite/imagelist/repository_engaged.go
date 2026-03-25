package imagelist

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapEngagedRows(rows []gen.ListImagesEngagedRow) ([]persist.ImageWithCursorKey[model.ScoreType], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesEngagedRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesEngagedRow) model.ScoreType {
			return model.ScoreType(row.ScoreEvaluated)
		},
	)
}

func mapEngagedAfterRows(rows []gen.ListImagesEngagedAfterRow) ([]persist.ImageWithCursorKey[model.ScoreType], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesEngagedAfterRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesEngagedAfterRow) model.ScoreType {
			return model.ScoreType(row.ScoreEvaluated)
		},
	)
}

func (repo repository) ListEngaged(
	ctx context.Context,
	spec persist.EngagedImageListSpec,
) ([]persist.ImageWithCursorKey[int16], error) {
	cursor, present := spec.CursorKey.Get()
	if !present {
		rows, err := repo.queries.ListImagesEngaged(
			ctx,
			gen.ListImagesEngagedParams{
				Limit:          int64(spec.MaxCount),
				ScoreThreshold: spec.ScoreThreshold,
			},
		)
		if err != nil {
			return nil, shared.MapSQLiteError("ListEngaged", err)
		}

		return mapEngagedRows(rows)
	}

	rows, err := repo.queries.ListImagesEngagedAfter(
		ctx,
		gen.ListImagesEngagedAfterParams{
			ScoreThreshold: spec.ScoreThreshold,
			CursorInt:      cursor.Primary,
			CursorID:       cursor.Secondary,
			MaxCount:       int64(spec.MaxCount),
		},
	)
	if err != nil {
		return nil, shared.MapSQLiteError("ListEngaged", err)
	}

	return mapEngagedAfterRows(rows)
}
