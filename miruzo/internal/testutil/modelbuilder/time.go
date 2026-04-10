package modelbuilder

import (
	"time"

	"github.com/samber/mo"
	"golang.org/x/exp/constraints"
)

func addSeconds[I constraints.Integer](baseTime time.Time, seconds I) time.Time {
	return baseTime.Add(time.Duration(seconds) * time.Second)
}

func resolveOffsetTime(v any, baseTime time.Time) mo.Option[time.Time] {
	switch value := v.(type) {
	case time.Duration:
		return mo.Some(baseTime.Add(value))
	case mo.Option[time.Duration]:
		if duration, present := value.Get(); present {
			return mo.Some(baseTime.Add(duration))
		}
		return mo.None[time.Time]()
	case int:
		return mo.Some(addSeconds(baseTime, value))
	case int8:
		return mo.Some(addSeconds(baseTime, value))
	case int16:
		return mo.Some(addSeconds(baseTime, value))
	case int32:
		return mo.Some(addSeconds(baseTime, value))
	case int64:
		return mo.Some(addSeconds(baseTime, value))
	case uint:
		return mo.Some(addSeconds(baseTime, value))
	case uint8:
		return mo.Some(addSeconds(baseTime, value))
	case uint16:
		return mo.Some(addSeconds(baseTime, value))
	case uint32:
		return mo.Some(addSeconds(baseTime, value))
	case uint64:
		return mo.Some(addSeconds(baseTime, value))
	case string:
		// Accept Go duration strings (e.g. "300ms", "1.5h", "-2m") via time.ParseDuration.
		if duration, err := time.ParseDuration(value); err == nil {
			return mo.Some(baseTime.Add(duration))
		}
		panic("invalid offset")
	default:
		panic("invalid offset")
	}
}
