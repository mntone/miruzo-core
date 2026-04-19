package backend

type Backend string

const (
	MySQL      Backend = "mysql"
	PostgreSQL Backend = "postgresql"
	SQLite     Backend = "sqlite"
)

func (b Backend) String() string {
	switch b {
	case MySQL:
		return "MySQL"
	case PostgreSQL:
		return "PostgreSQL"
	case SQLite:
		return "SQLite"
	default:
		return string(b)
	}
}
