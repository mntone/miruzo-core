package validate

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"golang.org/x/exp/constraints"
)

func ValidateRangeQuery[I constraints.Integer](key string, val, min, max I) *apierror.FieldError {
	if val < min || max < val {
		return &apierror.FieldError{
			Type:    apierror.FieldErrorTypeInvalid,
			Path:    "query." + key,
			Message: fmt.Sprintf("must be between %d and %d", min, max),
		}
	}

	return nil
}
