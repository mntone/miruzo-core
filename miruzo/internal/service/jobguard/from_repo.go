package jobguard

import (
	"context"
	"errors"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type jobRunDatabaseGuard struct {
	repository persist.JobRepository
}

func NewWithJobRepository(repository persist.JobRepository) JobRunGuard {
	return jobRunDatabaseGuard{
		repository: repository,
	}
}

func (grd jobRunDatabaseGuard) TryAcquire(
	ctx context.Context,
	name string,
	startedAt time.Time,
) (acquired bool, err error) {
	err = grd.repository.MarkStarted(ctx, name, startedAt)
	if err != nil {
		if errors.Is(err, persist.ErrConflict) {
			err = nil
		}
		return
	}

	acquired = true
	return
}

func (grd jobRunDatabaseGuard) Release(
	ctx context.Context,
	name string,
	finishedAt time.Time,
) error {
	err := grd.repository.MarkFinished(ctx, name, finishedAt)
	if err != nil && !errors.Is(err, persist.ErrConflict) {
		return err
	}

	return nil
}
