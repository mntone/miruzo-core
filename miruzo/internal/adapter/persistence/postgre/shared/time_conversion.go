package shared

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/mo"
)

var errOptionNoSuchElement = errors.New("no such element")

func TimeFromPgtype(value pgtype.Timestamp) time.Time {
	if !value.Valid {
		panic(errOptionNoSuchElement)
	}

	return value.Time
}

func OptionTimeFromPgtype(value pgtype.Timestamp) mo.Option[time.Time] {
	if !value.Valid {
		return mo.None[time.Time]()
	}
	return mo.Some(value.Time)
}

func PgtypeTimestampFromTime(value time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  value,
		Valid: true,
	}
}

func PgtypeTimestampFromOption(value mo.Option[time.Time]) pgtype.Timestamp {
	time, present := value.Get()
	if !present {
		return pgtype.Timestamp{}
	}
	return pgtype.Timestamp{
		Time:  time,
		Valid: true,
	}
}
