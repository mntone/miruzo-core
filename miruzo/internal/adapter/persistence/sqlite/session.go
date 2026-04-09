package sqlite

import (
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/action"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/imagelist"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/stats"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/user"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type sqliteRepositories struct {
	queries *gen.Queries
}

func (repos sqliteRepositories) ImageList() persist.ImageListRepository {
	return imagelist.NewRepository(repos.queries)
}

func (repos sqliteRepositories) Job() persist.JobRepository {
	return NewJobRepository(repos.queries)
}

func (repos sqliteRepositories) Settings() persist.SettingsRepository {
	return NewSettingsRepository(repos.queries)
}

func (repos sqliteRepositories) User() persist.UserRepository {
	return user.NewRepository(repos.queries)
}

type sqliteSessionRepositories struct {
	sqliteRepositories
}

func (repos sqliteSessionRepositories) Action() persist.ActionRepository {
	return action.NewRepository(repos.queries)
}

func (repos sqliteSessionRepositories) Stats() persist.StatsRepository {
	return stats.NewRepository(repos.queries)
}

func (repos sqliteSessionRepositories) View() persist.ViewRepository {
	return NewViewRepository(repos.queries)
}

func (repos sqliteSessionRepositories) User() persist.SessionUserRepository {
	return user.NewRepository(repos.queries)
}
