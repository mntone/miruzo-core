package shared

import (
	"database/sql"
	"errors"
	"time"

	"github.com/samber/mo"
)

var errOptionNoSuchElement = errors.New("no such element")

func TimeFromSql(value sql.NullTime) time.Time {
	if !value.Valid {
		panic(errOptionNoSuchElement)
	}

	return value.Time
}

func NullTimeFromTime(value time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  value,
		Valid: true,
	}
}

func NullTimeFromOption(value mo.Option[time.Time]) sql.NullTime {
	time, present := value.Get()
	if !present {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  time,
		Valid: true,
	}
}
