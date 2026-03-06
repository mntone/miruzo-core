package bind

import (
	"net/url"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"github.com/samber/mo"
)

const iso8601UTCLayout = "2006-01-02T15:04:05.999999Z"

func newInvalidTimeQueryErrors(queryName string) []apierror.FieldError {
	return []apierror.FieldError{{
		Path:    "query." + queryName,
		Type:    "invalid_type",
		Message: queryName + " must be a UTC timestamp in the format " + iso8601UTCLayout,
	}}
}

func ParseTimeQuery(
	queryValues url.Values,
	queryName string,
) (mo.Option[time.Time], []apierror.FieldError) {
	text := queryValues.Get(queryName)
	if text == "" {
		return mo.None[time.Time](), nil
	}

	parsedValue, parseError := time.Parse(iso8601UTCLayout, text)
	if parseError != nil {
		return mo.None[time.Time](), newInvalidTimeQueryErrors(queryName)
	}

	return mo.Some(parsedValue), nil
}

func ParseTimeQueryWithDefault(
	queryValues url.Values,
	queryName string,
	defaultValue time.Time,
) (time.Time, []apierror.FieldError) {
	text := queryValues.Get(queryName)
	if text == "" {
		return defaultValue, nil
	}

	parsedValue, parseError := time.Parse(iso8601UTCLayout, text)
	if parseError != nil {
		return time.Time{}, newInvalidTimeQueryErrors(queryName)
	}

	return parsedValue, nil
}
