package bind

import (
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

func BindUintQuery[U constraints.Unsigned](
	key string,
	values []string,
) (U, *apierror.FieldError) {
	text, fieldError := validateSingleValue(key, values)
	if fieldError != nil {
		return 0, fieldError
	}

	bitSize := bitSizeOfUnsignedInteger[U]()
	parsedValue64, parseError := strconv.ParseUint(text, 10, bitSize)
	if parseError != nil {
		return 0, &apierror.FieldError{
			Type:    apierror.FieldErrorTypeInvalid,
			Path:    "query." + key,
			Message: "must be an integer",
		}
	}

	return U(parsedValue64), nil
}
