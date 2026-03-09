package persist

type RepositoryFactory interface {
	Close() error

	NewAction() ActionRepository
	NewImageList() ImageListRepository
	NewUser() UserRepository
}
