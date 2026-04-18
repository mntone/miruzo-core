package mysql

import (
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/mysql/imagelist"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/mysql/stats"
	"github.com/mntone/miruzo-core/miruzo/internal/database/mysql/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type mysqlRepositories struct {
	queries *gen.Queries
}

func newRepositories(queries *gen.Queries) mysqlRepositories {
	return mysqlRepositories{
		queries: queries,
	}
}

func (repos mysqlRepositories) ImageList() persist.ImageListRepository {
	return imagelist.NewRepository(repos.queries)
}

func (repos mysqlRepositories) Job() persist.JobRepository {
	return jobRepository{
		queries: repos.queries,
	}
}

func (repos mysqlRepositories) Settings() persist.SettingsRepository {
	return settingsRepository{
		queries: repos.queries,
	}
}

func (repos mysqlRepositories) User() persist.UserRepository {
	return userRepository{
		queries: repos.queries,
	}
}

type mysqlSessionRepositories struct {
	mysqlRepositories
}

func NewSessionRepositories(queries *gen.Queries) mysqlSessionRepositories {
	return mysqlSessionRepositories{
		mysqlRepositories: newRepositories(queries),
	}
}

func (repos mysqlSessionRepositories) Action() persist.ActionRepository {
	return actionRepository{
		queries: repos.queries,
	}
}

func (repos mysqlSessionRepositories) Stats() persist.StatsRepository {
	return stats.NewRepository(repos.queries)
}

func (repos mysqlSessionRepositories) View() persist.ViewRepository {
	return nil
}

func (repos mysqlSessionRepositories) User() persist.SessionUserRepository {
	return userRepository{
		queries: repos.queries,
	}
}
