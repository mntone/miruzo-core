package bind

import (
	"net/url"
	"strconv"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"golang.org/x/exp/constraints"
)

func bitSizeOfUnsignedInteger[T constraints.Unsigned]() int {
	var zero T
	switch any(zero).(type) {
	case uint8:
		return 8
	case uint16:
		return 16
	case uint32:
		return 32
	case uint64:
		return 64
	case uint:
		return strconv.IntSize
	default:
		return 64
	}
}

func ParseUintQueryWithDefault[T constraints.Unsigned](
	queryValues url.Values,
	queryName string,
	defaultValue T,
) (T, []apierror.FieldError) {
	text := queryValues.Get(queryName)
	if text == "" {
		return defaultValue, nil
	}

	bitSize := bitSizeOfUnsignedInteger[T]()
	parsedValue64, parseError := strconv.ParseUint(text, 10, bitSize)
	if parseError != nil {
		return 0, []apierror.FieldError{{
			Path:    "query." + queryName,
			Type:    "invalid_type",
			Message: queryName + " must be an integer",
		}}
	}

	return T(parsedValue64), nil
}
