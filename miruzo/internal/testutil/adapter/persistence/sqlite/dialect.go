package sqlite

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
)

type sqliteDialect struct{}

func (sqliteDialect) Backend() backend.Backend {
	return backend.SQLite
}

func (sqliteDialect) ToPersistError(
	operation string,
	err error,
	mapping contract.DBErrorMapping,
) error {
	switch mapping {
	case contract.DBErrorMappingDefault:
		err = dberrors.ToPersist(operation, err)
	case contract.DBErrorMappingDelete:
		err = dberrors.ToPersistDelete(operation, err)
	}
	return err
}

func (sqliteDialect) BindVarStyle() contract.BindVarStyle {
	return contract.BindVarStyleQuestion
}

func (sqliteDialect) Param(index int32) string {
	return fmt.Sprintf("?%d", index)
}

func (d sqliteDialect) ParamRange(start, end int32) []any {
	return contract.ParamRange(start, end, d.Param)
}
