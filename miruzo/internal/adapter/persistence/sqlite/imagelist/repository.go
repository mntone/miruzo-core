package imagelist

import (
	"database/sql"

	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
)

type repository struct {
	queries *gen.Queries
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		queries: gen.New(db),
	}
}
