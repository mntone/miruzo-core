package context

import (
	"net/url"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/bind"
	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
)

const (
	levelDefault = "default"
	levelRich    = "rich"
)

type request struct {
	IsRich         bool
	ExcludeFormats []media.ImageFormat
}

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

func bindParams(queryValues url.Values) (request, []apierror.FieldError) {
	params := request{}

	var errors apierror.FieldErrors
	for key, values := range queryValues {
		switch key {
		case "level":
			ret, err := bindLevelQuery(key, values)
			if err != nil {
				errors = append(errors, *err)
				continue
			}

			params.IsRich = ret

		case "exclude_formats":
			excludeFormats, err := variant.BindImageFormatsQuery(key, values)
			if err != nil {
				errors = append(errors, *err)
				continue
			}

			params.ExcludeFormats = excludeFormats

		default:
			errors = append(errors, apierror.NewUnsupportedError(key))
		}
	}
	if errors != nil {
		errors.Sort()
	}

	return params, errors
}
