package sqlite

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	dbshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
)

type sqliteDialect struct{}

func (s sqliteDialect) Backend() backend.Backend {
	return backend.SQLite
}

func (s sqliteDialect) MapError(
	operation string,
	err error,
	mapping contract.DBErrorMapping,
) error {
	switch mapping {
	case contract.DBErrorMappingDefault:
		err = dbshared.MapSQLiteError(operation, err)
	case contract.DBErrorMappingDelete:
		err = dbshared.MapSQLiteDeleteError(operation, err)
	}
	return err
}

func (s sqliteDialect) BindVarStyle() contract.BindVarStyle {
	return contract.BindVarStyleQuestion
}

func (s sqliteDialect) Param(index int32) string {
	return fmt.Sprintf("?%d", index)
}

func (s sqliteDialect) ParamRange(start, end int32) []any {
	return contract.ParamRange(start, end, s.Param)
}
