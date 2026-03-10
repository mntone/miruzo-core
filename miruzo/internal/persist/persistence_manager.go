package persist

type Repositories struct {
	Action    ActionRepository
	ImageList ImageListRepository
	User      UserRepository
	View      ViewRepository
}

type PersistenceManager interface {
	Close() error

	Repos() Repositories
}
