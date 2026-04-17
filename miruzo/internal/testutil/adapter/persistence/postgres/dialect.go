package postgres

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
)

type postgresDialect struct{}

func (postgresDialect) Backend() backend.Backend {
	return backend.PostgreSQL
}

func (postgresDialect) MapError(
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

func (postgresDialect) BindVarStyle() contract.BindVarStyle {
	return contract.BindVarStyleDollar
}

func (postgresDialect) Param(index int32) string {
	return fmt.Sprintf("$%d", index)
}

func (d postgresDialect) ParamRange(start, end int32) []any {
	return contract.ParamRange(start, end, d.Param)
}
