package imagelist

import (
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
)

type repository struct {
	queries *gen.Queries
}

func NewRepository(queries *gen.Queries) repository {
	return repository{
		queries: queries,
	}
}
