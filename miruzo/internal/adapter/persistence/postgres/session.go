package postgres

import (
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/imagelist"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/stats"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type postgresRepositories struct {
	queries *gen.Queries
}

func newRepositories(queries *gen.Queries) postgresRepositories {
	return postgresRepositories{
		queries: queries,
	}
}

func (repos postgresRepositories) ImageList() persist.ImageListRepository {
	return imagelist.NewRepository(repos.queries)
}

func (repos postgresRepositories) Job() persist.JobRepository {
	return jobRepository{
		queries: repos.queries,
	}
}

func (repos postgresRepositories) Settings() persist.SettingsRepository {
	return settingsRepository{
		queries: repos.queries,
	}
}

func (repos postgresRepositories) User() persist.UserRepository {
	return userRepository{
		queries: repos.queries,
	}
}

type postgresSessionRepositories struct {
	postgresRepositories
}

func NewSessionRepositories(queries *gen.Queries) postgresSessionRepositories {
	return postgresSessionRepositories{
		postgresRepositories: newRepositories(queries),
	}
}

func (repos postgresSessionRepositories) Action() persist.ActionRepository {
	return actionRepository{
		queries: repos.queries,
	}
}

func (repos postgresSessionRepositories) Stats() persist.StatsRepository {
	return stats.NewRepository(repos.queries)
}

func (repos postgresSessionRepositories) User() persist.SessionUserRepository {
	return userRepository{
		queries: repos.queries,
	}
}

func (repos postgresSessionRepositories) View() persist.ViewRepository {
	return viewRepository{
		queries: repos.queries,
	}
}
