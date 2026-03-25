package list

import (
	"net/url"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/bind"
	"github.com/mntone/miruzo-core/miruzo/internal/api/validate"
	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	service "github.com/mntone/miruzo-core/miruzo/internal/service/imagelist"
	"github.com/samber/mo"
)

func bindParamsOf[S model.ImageListCursorScalar](
	queryValues url.Values,
	bindCursor func(key string, values []string) (model.ImageListCursorKey[S], *apierror.FieldError),
) (*service.Params[S], []apierror.FieldError) {
	params := service.Params[S]{
		Limit: defaultLimit,
	}

	var errors apierror.FieldErrors
	for key, values := range queryValues {
		switch key {
		case "limit":
			limit, err := bind.BindUintQuery[uint16](key, values)
			if err != nil {
				errors = append(errors, *err)
				continue
			}

			err = validate.ValidateRangeQuery(key, limit, limitMinimum, limitMaximum)
			if err != nil {
				errors = append(errors, *err)
				continue
			}

			params.Limit = limit

		case "cursor":
			cursor, err := bindCursor(key, values)
			if err != nil {
				errors = append(errors, *err)
				continue
			}

			params.Cursor = mo.Some(cursor)

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

	return &params, errors
}

func bindTimeCursorQuery(
	key string,
	values []string,
	expectedMode imageListCursorMode,
) (model.ImageListCursorKey[time.Time], *apierror.FieldError) {
	text, err := bind.ValidateSingleValue(key, values)
	if err != nil {
		return model.ImageListCursorKey[time.Time]{}, err
	}

	cursor, decodeError := decodeTimeImageListCursor(text, expectedMode)
	if decodeError != nil {
		return model.ImageListCursorKey[time.Time]{}, &apierror.FieldError{
			Type:    apierror.FieldErrorTypeInvalid,
			Path:    "query." + key,
			Message: "must be a valid cursor",
		}
	}
	return cursor, nil
}

func bindParamsForTimeCursor(
	queryValues url.Values,
	expectedMode imageListCursorMode,
) (*service.Params[time.Time], []apierror.FieldError) {
	return bindParamsOf(
		queryValues,
		func(key string, values []string) (model.ImageListCursorKey[time.Time], *apierror.FieldError) {
			return bindTimeCursorQuery(key, values, expectedMode)
		},
	)
}

func bindUint8CursorQuery(
	key string,
	values []string,
	expectedMode imageListCursorMode,
) (model.ImageListCursorKey[model.ScoreType], *apierror.FieldError) {
	text, err := bind.ValidateSingleValue(key, values)
	if err != nil {
		return model.ImageListCursorKey[model.ScoreType]{}, err
	}

	cursor, decodeError := decodeUint8ImageListCursor(text, expectedMode)
	if decodeError != nil {
		return model.ImageListCursorKey[model.ScoreType]{}, &apierror.FieldError{
			Type:    apierror.FieldErrorTypeInvalid,
			Path:    "query." + key,
			Message: "must be a valid cursor",
		}
	}
	return cursor, nil
}

func bindParamsForScoreCursor(
	queryValues url.Values,
	expectedMode imageListCursorMode,
) (*service.Params[model.ScoreType], []apierror.FieldError) {
	return bindParamsOf(
		queryValues,
		func(key string, values []string) (model.ImageListCursorKey[model.ScoreType], *apierror.FieldError) {
			return bindUint8CursorQuery(key, values, expectedMode)
		},
	)
}
