package stub

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type keyValuePair struct {
	Key   string
	Value string
}

type SettingsRepository struct {
	Store map[string]string

	GetError    error
	UpdateError error
	Updates     []keyValuePair
}

func NewStubSettingsRepository() *SettingsRepository {
	return &SettingsRepository{
		Store: make(map[string]string),
	}
}

func NewStubSettingsRepositoryWithKeyValue(
	key string,
	value string,
) *SettingsRepository {
	store := make(map[string]string)
	store[key] = value
	return &SettingsRepository{
		Store: store,
	}
}

func NewStubSettingsRepositoryWithStore(store map[string]string) *SettingsRepository {
	return &SettingsRepository{
		Store: store,
	}
}

func NewStubSettingsRepositoryWithGetError(getError error) *SettingsRepository {
	return &SettingsRepository{
		Store:    make(map[string]string),
		GetError: getError,
	}
}

func (repo *SettingsRepository) GetValue(
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

func (repo *SettingsRepository) UpdateValue(
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
