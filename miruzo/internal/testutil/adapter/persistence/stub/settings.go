package stub

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type keyValuePair struct {
	Key   string
	Value string
}

type settingsRepository struct {
	Store map[string]string

	GetError    error
	UpdateError error
	Updates     []keyValuePair
}

func NewStubSettingsRepository() *settingsRepository {
	return &settingsRepository{
		Store: make(map[string]string),
	}
}

func NewStubSettingsRepositoryWithKeyValue(
	key string,
	value string,
) *settingsRepository {
	store := make(map[string]string)
	store[key] = value
	return &settingsRepository{
		Store: store,
	}
}

func NewStubSettingsRepositoryWithStore(store map[string]string) *settingsRepository {
	return &settingsRepository{
		Store: store,
	}
}

func NewStubSettingsRepositoryWithGetError(getError error) *settingsRepository {
	return &settingsRepository{
		Store:    make(map[string]string),
		GetError: getError,
	}
}

func (repo *settingsRepository) GetValue(
	_ context.Context,
	key string,
) (string, error) {
	if repo.GetError != nil {
		return "", repo.GetError
	}

	value, ok := repo.Store[key]
	if !ok {
		return "", persist.ErrNotFound
	}

	return value, nil
}

func (repo *settingsRepository) UpdateValue(
	_ context.Context,
	key string,
	value string,
) error {
	repo.Updates = append(repo.Updates, keyValuePair{
		Key:   key,
		Value: value,
	})
	repo.Store[key] = value
	return repo.UpdateError
}
