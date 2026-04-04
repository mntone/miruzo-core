package postgres

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type jobRepository struct {
	queries *gen.Queries
}

func NewJobRepository(queries *gen.Queries) jobRepository {
	return jobRepository{
		queries: queries,
	}
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
		return shared.MapPostgreError("MarkStarted", err)
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
		FinishedAt: &finishedAt,
	})
	if err != nil {
		return shared.MapPostgreError("MarkFinished", err)
	}

	if rowCount == 0 {
		return persist.ErrConflict
	}

	return nil
}
