package list

import (
	"net/http"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/httperror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/response"
)

func (hdl *handler) listEngaged(
	responseWriter http.ResponseWriter,
	req *http.Request,
) {
	params, fieldError := bindParamsForScoreCursor(req.URL.Query())
	if fieldError != nil {
		response.WriteJSON(
			responseWriter,
			http.StatusBadRequest,
			apierror.NewValidationError(fieldError),
		)
		return
	}

	result, serviceError := hdl.service.ListEngaged(req.Context(), params)
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
