package persist

import "context"

type Repositories struct {
	Action    ActionRepository
	ImageList ImageListRepository
	Settings  SettingsRepository
	Stats     StatsRepository
	StatsList StatsListRepository
	User      UserRepository
	View      ViewRepository
}

type SessionCallback func(ctx context.Context, repos Repositories) error

type PersistenceManager interface {
	Close() error

	Repos() Repositories
	Session(
		ctx context.Context,
		callback SessionCallback,
	) error
}
