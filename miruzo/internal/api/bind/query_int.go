package bind

import (
	"net/url"
	"strconv"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"github.com/samber/mo"
	"golang.org/x/exp/constraints"
)

func ParseIntQuery[T constraints.Signed](
	queryValues url.Values,
	queryName string,
) (mo.Option[T], []apierror.FieldError) {
	text := queryValues.Get(queryName)
	if text == "" {
		return mo.None[T](), nil
	}

	bitSize := bitSizeOfSignedInteger[T]()
	parsedValue64, parseError := strconv.ParseInt(text, 10, bitSize)
	if parseError != nil {
		return mo.None[T](), []apierror.FieldError{{
			Path:    "query." + queryName,
			Type:    "invalid_type",
			Message: queryName + " must be an integer",
		}}
	}

	return mo.Some(T(parsedValue64)), nil
}

func ParseIntQueryWithDefault[T constraints.Signed](
	queryValues url.Values,
	queryName string,
	defaultValue T,
) (T, []apierror.FieldError) {
	text := queryValues.Get(queryName)
	if text == "" {
		return defaultValue, nil
	}

	bitSize := bitSizeOfSignedInteger[T]()
	parsedValue64, parseError := strconv.ParseInt(text, 10, bitSize)
	if parseError != nil {
		return 0, []apierror.FieldError{{
			Path:    "query." + queryName,
			Type:    "invalid_type",
			Message: queryName + " must be an integer",
		}}
	}

	return T(parsedValue64), nil
}
