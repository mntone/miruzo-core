package backend

type Backend string

const (
	MySQL      Backend = "mysql"
	PostgreSQL Backend = "postgres"
	SQLite     Backend = "sqlite"
)
