package persist

type Repositories struct {
	Action    ActionRepository
	ImageList ImageListRepository
	User      UserRepository
}

type PersistenceManager interface {
	Close() error

	Repos() Repositories
}
