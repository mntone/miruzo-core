package bind

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
)

const iso8601UTCLayout = "2006-01-02T15:04:05.999999Z"

func BindTimeQuery(
	key string,
	values []string,
) (time.Time, *apierror.FieldError) {
	text, fieldError := ValidateSingleValue(key, values)
	if fieldError != nil {
		return time.Time{}, fieldError
	}

	parsedValue, parseError := time.Parse(iso8601UTCLayout, text)
	if parseError != nil {
		return time.Time{}, &apierror.FieldError{
			Type:    apierror.FieldErrorTypeInvalid,
			Path:    "query." + key,
			Message: "must be ISO8601 timestamp",
		}
	}

	return parsedValue, nil
}
