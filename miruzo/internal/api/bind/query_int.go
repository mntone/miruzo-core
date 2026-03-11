package bind

import (
	"strconv"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"golang.org/x/exp/constraints"
)

func BindIntQuery[I constraints.Signed](
	key string,
	values []string,
) (I, *apierror.FieldError) {
	text, fieldError := validateSingleValue(key, values)
	if fieldError != nil {
		return 0, fieldError
	}

	bitSize := bitSizeOfSignedInteger[I]()
	parsedValue64, parseError := strconv.ParseInt(text, 10, bitSize)
	if parseError != nil {
		return 0, &apierror.FieldError{
			Type:    apierror.FieldErrorTypeInvalid,
			Path:    "query." + key,
			Message: "must be an integer",
		}
	}

	return I(parsedValue64), nil
}
