package stub

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type actionRepositoryCreateArgs struct {
	IngestID      model.IngestIDType
	Type          model.ActionType
	OccurredAt    time.Time
	PeriodStartAt time.Time
}

type actionRepositoryExistsSinceArgs struct {
	IngestID        model.IngestIDType
	Type            model.ActionType
	SinceOccurredAt time.Time
}

type actionStorage struct {
	Store  []model.Action
	NextID model.ActionIDType
}

type actionRepository struct {
	actionStorage

	CreateError      error
	CreateArgs       []actionRepositoryCreateArgs
	ExistsSinceError error
	ExistsSinceArgs  []actionRepositoryExistsSinceArgs
}

func NewStubActionRepository() *actionRepository {
	return &actionRepository{
		actionStorage: actionStorage{
			NextID: 1,
		},
	}
}

func NewStubActionRepositoryWithNextID(nextID model.ActionIDType) *actionRepository {
	return &actionRepository{
		actionStorage: actionStorage{
			NextID: nextID,
		},
	}
}

func (repo actionRepository) snapshot() actionStorage {
	var store []model.Action
	if repo.Store != nil {
		store = make([]model.Action, len(repo.Store))
		copy(store, repo.Store)
	}
	return actionStorage{
		Store:  store,
		NextID: repo.NextID,
	}
}

func (repo *actionRepository) Create(
	_ context.Context,
	ingestID model.IngestIDType,
	kind model.ActionType,
	occurredAt time.Time,
	periodStartAt time.Time,
) (model.ActionIDType, error) {
	repo.CreateArgs = append(repo.CreateArgs, actionRepositoryCreateArgs{
		IngestID:      ingestID,
		Type:          kind,
		OccurredAt:    occurredAt,
		PeriodStartAt: periodStartAt,
	})

	if repo.CreateError != nil {
		return 0, repo.CreateError
	}

	action := model.Action{
		ID:         repo.NextID,
		IngestID:   ingestID,
		Type:       kind,
		OccurredAt: occurredAt,
	}
	repo.Store = append(repo.Store, action)

	repo.NextID += 1
	return action.ID, nil
}

func (repo *actionRepository) ExistsSince(
	_ context.Context,
	ingestID model.IngestIDType,
	kind model.ActionType,
	sinceOccurredAt time.Time,
) (bool, error) {
	repo.ExistsSinceArgs = append(repo.ExistsSinceArgs, actionRepositoryExistsSinceArgs{
		IngestID:        ingestID,
		Type:            kind,
		SinceOccurredAt: sinceOccurredAt,
	})

	if repo.ExistsSinceError != nil {
		return false, repo.ExistsSinceError
	}

	for _, action := range repo.Store {
		if action.IngestID != ingestID {
			continue
		}

		if action.Type != kind {
			continue
		}

		if action.OccurredAt.Before(sinceOccurredAt) {
			continue
		}

		return true, nil
	}

	return false, nil
}
