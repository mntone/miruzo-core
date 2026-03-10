package persist

type Repositories struct {
	Action    ActionRepository
	ImageList ImageListRepository
	Stats     StatsRepository
	User      UserRepository
	View      ViewRepository
}

type PersistenceManager interface {
	Close() error

	Repos() Repositories
}
