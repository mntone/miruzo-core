package postgres

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

const (
	jobStartedOperationName  = "MarkStarted"
	jobFinishedOperationName = "MarkFinished"
)

type jobRepository struct {
	queries *gen.Queries
}

func (repo jobRepository) MarkStarted(
	ctx context.Context,
	name string,
	startedAt time.Time,
) error {
	affectedRows, err := repo.queries.MarkJobStarted(ctx, gen.MarkJobStartedParams{
		Name:      name,
		StartedAt: startedAt,
	})
	if err != nil {
		return dberrors.ToPersist(jobStartedOperationName, err)
	}
	if affectedRows == 1 {
		return nil
	}

	var baseError error
	if affectedRows == 0 {
		baseError = persist.ErrConflict
	} else {
		baseError = persist.ErrInvariantViolation
	}
	return dberrors.WrapKV(
		baseError,
		jobStartedOperationName,
		"affected_rows", affectedRows,
		"name", name,
	)
}

func (repo jobRepository) MarkFinished(
	ctx context.Context,
	name string,
	finishedAt time.Time,
) error {
	affectedRows, err := repo.queries.MarkJobFinished(ctx, gen.MarkJobFinishedParams{
		Name:       name,
		FinishedAt: &finishedAt,
	})
	if err != nil {
		return dberrors.ToPersist(jobFinishedOperationName, err)
	}
	if affectedRows == 1 {
		return nil
	}

	var baseError error
	if affectedRows == 0 {
		baseError = persist.ErrConflict
	} else {
		baseError = persist.ErrInvariantViolation
	}
	return dberrors.WrapKV(
		baseError,
		jobFinishedOperationName,
		"affected_rows", affectedRows,
		"name", name,
	)
}
