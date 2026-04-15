package stub

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type ActionModel struct {
	ID            model.ActionIDType
	IngestID      model.IngestIDType
	Type          model.ActionType
	OccurredAt    time.Time
	PeriodStartAt time.Time
}

type actionRepositoryCreateArgs struct {
	IngestID      model.IngestIDType
	Type          model.ActionType
	OccurredAt    time.Time
	PeriodStartAt time.Time
}

type actionRepositoryCreateDailyDecayIfAbsentArgs struct {
	IngestID      model.IngestIDType
	OccurredAt    time.Time
	PeriodStartAt time.Time
}

type actionRepositoryCreateLoveIfAbsentArgs struct {
	IngestID      model.IngestIDType
	Type          persist.LoveActionType
	OccurredAt    time.Time
	PeriodStartAt time.Time
}

type actionRepositoryExistsSinceArgs struct {
	IngestID        model.IngestIDType
	Type            model.ActionType
	SinceOccurredAt time.Time
}

type actionStorage struct {
	Store  []ActionModel
	NextID model.ActionIDType
}

type actionRepository struct {
	actionStorage

	CreateError                   error
	CreateArgs                    []actionRepositoryCreateArgs
	CreateDailyDecayIfAbsentError error
	CreateDailyDecayIfAbsentArgs  []actionRepositoryCreateDailyDecayIfAbsentArgs
	CreateLoveIfAbsentError       error
	CreateLoveIfAbsentArgs        []actionRepositoryCreateLoveIfAbsentArgs
	ExistsSinceError              error
	ExistsSinceArgs               []actionRepositoryExistsSinceArgs
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
	var store []ActionModel
	if repo.Store != nil {
		store = make([]ActionModel, len(repo.Store))
		copy(store, repo.Store)
	}
	return actionStorage{
		Store:  store,
		NextID: repo.NextID,
	}
}

func (repo *actionRepository) appendCreatedAction(
	ingestID model.IngestIDType,
	kind model.ActionType,
	occurredAt time.Time,
	periodStartAt time.Time,
) model.ActionIDType {
	actionID := repo.NextID
	repo.Store = append(repo.Store, ActionModel{
		ID:            repo.NextID,
		IngestID:      ingestID,
		Type:          kind,
		OccurredAt:    occurredAt,
		PeriodStartAt: periodStartAt,
	})

	repo.NextID += 1
	return actionID
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

	actionID := repo.appendCreatedAction(ingestID, kind, occurredAt, periodStartAt)
	return actionID, nil
}

func (repo *actionRepository) CreateDailyDecayIfAbsent(
	_ context.Context,
	ingestID model.IngestIDType,
	occurredAt time.Time,
	periodStartAt time.Time,
) error {
	repo.CreateDailyDecayIfAbsentArgs = append(repo.CreateDailyDecayIfAbsentArgs, actionRepositoryCreateDailyDecayIfAbsentArgs{
		IngestID:      ingestID,
		OccurredAt:    occurredAt,
		PeriodStartAt: periodStartAt,
	})

	if repo.CreateDailyDecayIfAbsentError != nil {
		return repo.CreateDailyDecayIfAbsentError
	}

	for _, action := range repo.Store {
		if action.IngestID != ingestID {
			continue
		}
		if action.Type != model.ActionTypeDecay {
			continue
		}
		if !action.PeriodStartAt.Equal(periodStartAt) {
			continue
		}
		return persist.ErrConflict
	}

	repo.appendCreatedAction(ingestID, model.ActionTypeDecay, occurredAt, periodStartAt)
	return nil
}

func (repo *actionRepository) CreateLoveIfAbsent(
	_ context.Context,
	ingestID model.IngestIDType,
	loveType persist.LoveActionType,
	occurredAt time.Time,
	periodStartAt time.Time,
) error {
	repo.CreateLoveIfAbsentArgs = append(repo.CreateLoveIfAbsentArgs, actionRepositoryCreateLoveIfAbsentArgs{
		IngestID:      ingestID,
		Type:          loveType,
		OccurredAt:    occurredAt,
		PeriodStartAt: periodStartAt,
	})

	if repo.CreateLoveIfAbsentError != nil {
		return repo.CreateLoveIfAbsentError
	}

	for _, action := range repo.Store {
		if action.IngestID != ingestID {
			continue
		}
		if action.Type != model.ActionType(persist.LoveActionTypeLove) &&
			action.Type != model.ActionType(persist.LoveActionTypeLoveCanceled) {
			continue
		}
		if !action.OccurredAt.Equal(occurredAt) {
			continue
		}
		return persist.ErrConflict
	}

	repo.appendCreatedAction(
		ingestID,
		model.ActionType(loveType),
		occurredAt,
		periodStartAt,
	)
	return nil
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
