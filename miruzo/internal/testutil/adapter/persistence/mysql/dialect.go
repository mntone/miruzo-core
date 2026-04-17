package mysql

import (
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	dbshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/mysql/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
)

type mysqlDialect struct{}

func (mysqlDialect) Backend() backend.Backend {
	return backend.MySQL
}

func (mysqlDialect) MapError(
	operation string,
	err error,
	mapping contract.DBErrorMapping,
) error {
	switch mapping {
	case contract.DBErrorMappingDefault, contract.DBErrorMappingDelete:
		err = dbshared.MapMySQLError(operation, err)
	}
	return err
}

func (mysqlDialect) BindVarStyle() contract.BindVarStyle {
	return contract.BindVarStyleQuestion
}

func (mysqlDialect) Param(index int32) string {
	return "?"
}

func (d mysqlDialect) ParamRange(start, end int32) []any {
	return contract.ParamRange(start, end, d.Param)
}
