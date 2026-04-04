package stub

import (
	"context"
	"maps"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

type jobRepositoryMarkStartedArgs struct {
	Name      string
	StartedAt time.Time
}

type jobRepositoryMarkFinishedArgs struct {
	Name       string
	FinishedAt time.Time
}

type jobStorage struct {
	Store map[string]model.Job
}

type jobRepository struct {
	jobStorage

	MarkStartedError  error
	MarkStartedArgs   []jobRepositoryMarkStartedArgs
	MarkFinishedError error
	MarkFinishedArgs  []jobRepositoryMarkFinishedArgs
}

func NewStubJobRepository(jobs ...model.Job) *jobRepository {
	store := make(map[string]model.Job, len(jobs))
	for _, j := range jobs {
		store[j.Name] = j
	}
	return &jobRepository{
		jobStorage: jobStorage{
			Store: store,
		},
	}
}

func (repo jobRepository) snapshot() jobStorage {
	var store map[string]model.Job
	if repo.Store != nil {
		store = make(map[string]model.Job, len(repo.Store))
		maps.Copy(store, repo.Store)
	}
	return jobStorage{
		Store: store,
	}
}

func (repo *jobRepository) MarkStarted(
	_ context.Context,
	name string,
	startedAt time.Time,
) error {
	repo.MarkStartedArgs = append(repo.MarkStartedArgs, jobRepositoryMarkStartedArgs{
		Name:      name,
		StartedAt: startedAt,
	})

	if repo.MarkStartedError != nil {
		return repo.MarkStartedError
	}

	job, ok := repo.Store[name]
	if !ok {
		job = model.Job{
			Name:      name,
			StartedAt: startedAt,
		}
	} else {
		if job.FinishedAt.IsAbsent() {
			return persist.ErrConflict
		}

		job.StartedAt = startedAt
		job.FinishedAt = mo.None[time.Time]()
	}

	repo.Store[name] = job
	return nil
}

func (repo *jobRepository) MarkFinished(
	_ context.Context,
	name string,
	finishedAt time.Time,
) error {
	repo.MarkFinishedArgs = append(repo.MarkFinishedArgs, jobRepositoryMarkFinishedArgs{
		Name:       name,
		FinishedAt: finishedAt,
	})

	if repo.MarkFinishedError != nil {
		return repo.MarkFinishedError
	}

	job, ok := repo.Store[name]
	if !ok {
		return persist.ErrConflict
	}
	if job.FinishedAt.IsPresent() {
		return persist.ErrConflict
	}

	job.FinishedAt = mo.Some(finishedAt)

	repo.Store[name] = job
	return nil
}
