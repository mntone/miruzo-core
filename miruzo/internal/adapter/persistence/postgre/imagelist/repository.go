package imagelist

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
)

type repository struct {
	queries *gen.Queries
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		queries: gen.New(pool),
	}
}
