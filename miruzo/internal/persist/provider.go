package persist

import "context"

type commonRepositories interface {
	ImageList() ImageListRepository
	Job() JobRepository
	Settings() SettingsRepository
}

type Repositories interface {
	commonRepositories
	User() UserRepository
}

type SessionRepositories interface {
	commonRepositories
	Action() ActionRepository
	Stats() StatsRepository
	User() SessionUserRepository
	View() ViewRepository
}

type SessionCallback func(ctx context.Context, repos SessionRepositories) error

type PersistenceProvider interface {
	Repos() Repositories
	Session(
		ctx context.Context,
		callback SessionCallback,
	) error
}
