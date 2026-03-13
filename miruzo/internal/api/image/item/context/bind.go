package context

import (
	"net/url"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/bind"
)

const (
	levelDefault = "default"
	levelRich    = "rich"
)

func bindLevelQuery(key string, values []string) (bool, *apierror.FieldError) {
	text, fieldError := bind.ValidateSingleValue(key, values)
	if fieldError != nil {
		return false, fieldError
	}

	switch text {
	case levelRich:
		return true, nil
	case levelDefault:
		return false, nil
	}
	return false, &apierror.FieldError{
		Type:    apierror.FieldErrorTypeInvalid,
		Path:    "query.level",
		Message: "must be \"default\" or \"rich\"",
	}
}

func bindParams(queryValues url.Values) (bool, []apierror.FieldError) {
	rich := false

	var errors apierror.FieldErrors
	for key, values := range queryValues {
		switch key {
		case "level":
			ret, err := bindLevelQuery(key, values)
			if err != nil {
				errors = append(errors, *err)
				continue
			}

			rich = ret

		default:
			errors = append(errors, apierror.NewUnsupportedError(key))
		}
	}
	if errors != nil {
		errors.Sort()
	}

	return rich, errors
}
