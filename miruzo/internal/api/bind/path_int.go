package bind

import (
	"net/http"
	"strconv"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"golang.org/x/exp/constraints"
)

func BindIntPath[S constraints.Signed](
	request *http.Request,
	pathName string,
) (S, *apierror.FieldError) {
	text := request.PathValue(pathName)
	if text == "" {
		return 0, &apierror.FieldError{
			Type:    apierror.FieldErrorTypeRequired,
			Path:    "path." + pathName,
			Message: "is required",
		}
	}

	bitSize := bitSizeOfSignedInteger[S]()
	parsedValue64, parseError := strconv.ParseInt(text, 10, bitSize)
	if parseError != nil {
		return 0, &apierror.FieldError{
			Type:    apierror.FieldErrorTypeInvalid,
			Path:    "path." + pathName,
			Message: "must be an integer",
		}
	}

	return S(parsedValue64), nil
}
