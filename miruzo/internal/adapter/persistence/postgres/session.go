package postgres

import (
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/action"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/imagelist"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/stats"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/user"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type postgresRepositories struct {
	queries *gen.Queries
}

func (repos postgresRepositories) ImageList() persist.ImageListRepository {
	return imagelist.NewRepository(repos.queries)
}

func (repos postgresRepositories) Job() persist.JobRepository {
	return NewJobRepository(repos.queries)
}

func (repos postgresRepositories) Settings() persist.SettingsRepository {
	return NewSettingsRepository(repos.queries)
}

func (repos postgresRepositories) User() persist.UserRepository {
	return user.NewRepository(repos.queries)
}

type postgresSessionRepositories struct {
	postgresRepositories
}

func (repos postgresSessionRepositories) Action() persist.ActionRepository {
	return action.NewRepository(repos.queries)
}

func (repos postgresSessionRepositories) Stats() persist.StatsRepository {
	return stats.NewRepository(repos.queries)
}

func (repos postgresSessionRepositories) User() persist.SessionUserRepository {
	return user.NewRepository(repos.queries)
}

func (repos postgresSessionRepositories) View() persist.ViewRepository {
	return NewViewRepository(repos.queries)
}
