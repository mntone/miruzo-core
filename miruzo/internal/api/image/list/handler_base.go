package list

import (
	"context"
	"net/http"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	httperror "github.com/mntone/miruzo-core/miruzo/internal/api/http/error"
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
	params, fieldError := bindParamsForTimeCursor(req.URL.Query())
	if fieldError != nil {
		response.WriteJSON(
			responseWriter,
			http.StatusBadRequest,
			apierror.NewValidationError(fieldError),
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
