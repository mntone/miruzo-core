package bind

import (
	"net/http"
	"strconv"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"golang.org/x/exp/constraints"
)

func ParseIntPath[T constraints.Signed](
	request *http.Request,
	pathName string,
) (T, []apierror.FieldError) {
	text := request.PathValue(pathName)
	if text == "" {
		return T(0), []apierror.FieldError{{
			Path:    "path." + pathName,
			Type:    "missing",
			Message: pathName + " is required",
		}}
	}

	bitSize := bitSizeOfSignedInteger[T]()
	parsedValue64, parseError := strconv.ParseInt(text, 10, bitSize)
	if parseError != nil {
		return T(0), []apierror.FieldError{{
			Path:    "path." + pathName,
			Type:    "invalid_type",
			Message: pathName + " must be an integer",
		}}
	}

	return T(parsedValue64), nil
}
