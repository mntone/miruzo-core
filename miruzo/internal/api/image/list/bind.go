package list

import (
	"net/url"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/bind"
	"github.com/mntone/miruzo-core/miruzo/internal/api/validate"
	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	service "github.com/mntone/miruzo-core/miruzo/internal/service/imagelist"
	"github.com/samber/mo"
)

func bindParamsOf[C persist.ImageListCursor](
	queryValues url.Values,
	bindCursor func(key string, values []string) (C, *apierror.FieldError),
) (*service.Params[C], []apierror.FieldError) {
	params := service.Params[C]{
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

func bindParamsForTimeCursor(queryValues url.Values) (*service.Params[time.Time], []apierror.FieldError) {
	return bindParamsOf(queryValues, bind.BindTimeQuery)
}

func bindParamsForScoreCursor(queryValues url.Values) (*service.Params[model.ScoreType], []apierror.FieldError) {
	return bindParamsOf(queryValues, bind.BindIntQuery[model.ScoreType])
}
