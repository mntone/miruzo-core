package postgres

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	dbshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
)

type postgresDialect struct{}

func (s postgresDialect) Backend() backend.Backend {
	return backend.PostgreSQL
}

func (s postgresDialect) MapError(
	operation string,
	err error,
	mapping contract.DBErrorMapping,
) error {
	switch mapping {
	case contract.DBErrorMappingDefault:
		err = dbshared.MapPostgreError(operation, err)
	case contract.DBErrorMappingDelete:
		err = dbshared.MapPostgreDeleteError(operation, err)
	}
	return err
}

func (s postgresDialect) BindVarStyle() contract.BindVarStyle {
	return contract.BindVarStyleDollar
}

func (s postgresDialect) Param(index int32) string {
	return fmt.Sprintf("$%d", index)
}

func (s postgresDialect) ParamRange(start, end int32) []any {
	return contract.ParamRange(start, end, s.Param)
}
