package list

import (
	"context"
	"net/http"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/httperror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/response"
	"github.com/mntone/miruzo-core/miruzo/internal/service/imagelist"
)

func (hdl *handler) listBase(
	responseWriter http.ResponseWriter,
	req *http.Request,
	listFn func(
		requestContext context.Context,
		params *imagelist.Params[time.Time],
	) (imagelist.Result[time.Time], error),
) {
	params, fieldErrors := buildTimeParamsFromQuery(req.URL.Query())
	if len(fieldErrors) != 0 {
		response.WriteJSON(
			responseWriter,
			http.StatusBadRequest,
			apierror.NewValidationError(fieldErrors),
		)
		return
	}

	result, serviceError := listFn(req.Context(), params)
	if serviceError != nil {
		httperror.WriteServiceError(responseWriter, serviceError)
		return
	}

	_ = response.WriteJSON(
		responseWriter,
		http.StatusOK,
		mapImageListResponse(
			result,
			hdl.variantLayersConfig,
			hdl.mediaURLBuilder,
		),
	)
}
