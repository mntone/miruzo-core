package list

import (
	"net/url"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/bind"
	"github.com/mntone/miruzo-core/miruzo/internal/api/common"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	imageListService "github.com/mntone/miruzo-core/miruzo/internal/service/imagelist"
)

type query[C persist.ImageListCursor] struct {
	common.PaginationQuery[C]
}

func bindTimeQueryWithTimeCursor(qry url.Values) (query[time.Time], []apierror.FieldError) {
	limit, errs := bind.ParseUintQueryWithDefault(qry, "limit", defaultLimit)
	if len(errs) != 0 {
		return query[time.Time]{}, errs
	}

	cursor, errs := bind.ParseTimeQuery(qry, "cursor")
	if len(errs) != 0 {
		return query[time.Time]{}, errs
	}

	return query[time.Time]{
		PaginationQuery: common.PaginationQuery[time.Time]{
			Limit:  limit,
			Cursor: cursor,
		},
	}, nil
}

func bindIntQueryWithTimeCursor(qry url.Values) (query[int16], []apierror.FieldError) {
	limit, errs := bind.ParseUintQueryWithDefault(qry, "limit", defaultLimit)
	if len(errs) != 0 {
		return query[int16]{}, errs
	}

	cursor, errs := bind.ParseIntQuery[int16](qry, "cursor")
	if len(errs) != 0 {
		return query[int16]{}, errs
	}

	return query[int16]{
		PaginationQuery: common.PaginationQuery[int16]{
			Limit:  limit,
			Cursor: cursor,
		},
	}, nil
}

func buildTimeParamsFromQuery(
	queryValues url.Values,
) (*imageListService.Params[time.Time], []apierror.FieldError) {
	qry, errs := bindTimeQueryWithTimeCursor(queryValues)
	if len(errs) != 0 {
		return nil, errs
	}

	errs = common.ValidatePaginationQuery(
		qry.PaginationQuery,
		limitMinimum,
		limitMaximum,
	)
	if len(errs) != 0 {
		return nil, errs
	}

	return &imageListService.Params[time.Time]{
		Cursor:         qry.Cursor,
		Limit:          qry.Limit,
		ExcludeFormats: nil,
	}, nil
}

func buildIntParamsFromQuery(
	queryValues url.Values,
) (*imageListService.Params[int16], []apierror.FieldError) {
	qry, errs := bindIntQueryWithTimeCursor(queryValues)
	if len(errs) != 0 {
		return nil, errs
	}

	errs = common.ValidatePaginationQuery(
		qry.PaginationQuery,
		limitMinimum,
		limitMaximum,
	)
	if len(errs) != 0 {
		return nil, errs
	}

	return &imageListService.Params[int16]{
		Cursor:         qry.Cursor,
		Limit:          qry.Limit,
		ExcludeFormats: nil,
	}, nil
}
