package stub

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type actionRepositoryCreateArgs struct {
	IngestID   model.IngestIDType
	Type       model.ActionType
	OccurredAt time.Time
}

type actionStorage struct {
	Store  []model.Action
	NextID model.ActionIDType
}

type actionRepository struct {
	actionStorage

	CreateError error
	CreateArgs  []actionRepositoryCreateArgs
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
) (model.ActionIDType, error) {
	repo.CreateArgs = append(repo.CreateArgs, actionRepositoryCreateArgs{
		IngestID:   ingestID,
		Type:       kind,
		OccurredAt: occurredAt,
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
