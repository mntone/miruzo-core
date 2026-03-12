package bind

import (
	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
)

func ValidateSingleValue(
	key string,
	values []string,
) (string, *apierror.FieldError) {
	l := len(values)
	if l != 1 {
		if l == 0 {
			return "", &apierror.FieldError{
				Type:    apierror.FieldErrorTypeInvalid,
				Path:    "query." + key,
				Message: "must not be empty",
			}
		}

		return "", &apierror.FieldError{
			Type:    apierror.FieldErrorTypeDuplicate,
			Path:    "query." + key,
			Message: "must not be specified multiple times",
		}
	}

	return values[0], nil
}
