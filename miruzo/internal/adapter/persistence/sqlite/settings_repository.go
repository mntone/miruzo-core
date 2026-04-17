package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type settingsRepository struct {
	queries *gen.Queries
}

func (repo settingsRepository) GetValue(
	ctx context.Context,
	key string,
) (string, error) {
	value, err := repo.queries.GetSettingsValueByKey(ctx, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", persist.ErrNoRows
		}

		return "", dberrors.ToPersist("GetValue", err)
	}

	return value, nil
}

func (repo settingsRepository) UpdateValue(
	ctx context.Context,
	key string,
	value string,
) error {
	err := repo.queries.UpdateSettingsValueByKey(ctx, gen.UpdateSettingsValueByKeyParams{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return dberrors.ToPersist("UpdateValue", err)
	}

	return nil
}
