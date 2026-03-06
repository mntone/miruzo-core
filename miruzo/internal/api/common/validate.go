package common

import "github.com/mntone/miruzo-core/miruzo/internal/api/apierror"

func ValidatePaginationQuery[C any](
	query PaginationQuery[C],
	limitMinimum uint16,
	limitMaximum uint16,
) []apierror.FieldError {
	fieldErrors := make([]apierror.FieldError, 0, 2)

	// limit
	if query.Limit < limitMinimum {
		fieldErrors = append(fieldErrors, apierror.FieldError{
			Path:    "query.limit",
			Type:    "out_of_range",
			Message: "limit is too small",
		})
	} else if query.Limit > limitMaximum {
		fieldErrors = append(fieldErrors, apierror.FieldError{
			Path:    "query.limit",
			Type:    "out_of_range",
			Message: "limit is too large",
		})
	}

	return fieldErrors
}
