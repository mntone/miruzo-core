package persist

type RepositoryFactory interface {
	Close() error

	NewImageList() ImageListRepository
	NewUser() UserRepository
}
