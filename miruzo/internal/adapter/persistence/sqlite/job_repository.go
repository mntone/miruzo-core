package sqlite

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type jobRepository struct {
	queries *gen.Queries
}

func (repo jobRepository) MarkStarted(
	ctx context.Context,
	name string,
	startedAt time.Time,
) error {
	rowCount, err := repo.queries.MarkJobStarted(ctx, gen.MarkJobStartedParams{
		Name:      name,
		StartedAt: startedAt,
	})
	if err != nil {
		return dberrors.ToPersist("MarkStarted", err)
	}

	if rowCount == 0 {
		return persist.ErrConflict
	}

	return nil
}

func (repo jobRepository) MarkFinished(
	ctx context.Context,
	name string,
	finishedAt time.Time,
) error {
	rowCount, err := repo.queries.MarkJobFinished(ctx, gen.MarkJobFinishedParams{
		Name:       name,
		FinishedAt: finishedAt,
	})
	if err != nil {
		return dberrors.ToPersist("MarkFinished", err)
	}

	if rowCount == 0 {
		return persist.ErrConflict
	}

	return nil
}
