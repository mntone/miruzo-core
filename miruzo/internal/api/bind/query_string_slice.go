package bind

import (
	"strings"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
)

func BindStringSliceQuery(
	key string,
	values []string,
	sep string,
) ([]string, *apierror.FieldError) {
	text, fieldError := ValidateSingleValue(key, values)
	if fieldError != nil {
		return nil, fieldError
	}

	parsedValue := make([]string, 0, strings.Count(text, sep))
	for rawValue := range strings.SplitSeq(text, sep) {
		if rawValue == "" {
			continue
		}

		parsedValue = append(parsedValue, rawValue)
	}
	if len(parsedValue) == 0 {
		return nil, &apierror.FieldError{
			Type:    apierror.FieldErrorTypeInvalid,
			Path:    "query." + key,
			Message: "must not be empty",
		}
	}

	return parsedValue, nil
}
