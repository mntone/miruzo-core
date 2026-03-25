package context

import (
	"net/http"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/bind"
	httperror "github.com/mntone/miruzo-core/miruzo/internal/api/http/error"
	"github.com/mntone/miruzo-core/miruzo/internal/api/response"
	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/service/view"
)

type handler struct {
	service         view.Service
	mediaURLBuilder variant.MediaURLBuilder
}

func NewHandler(
	srv view.Service,
	mediaURLBuilder variant.MediaURLBuilder,
) *handler {
	return &handler{
		service:         srv,
		mediaURLBuilder: mediaURLBuilder,
	}
}

func (hdl *handler) getContext(
	responseWriter http.ResponseWriter,
	req *http.Request,
) {
	ingestID, fieldError := bind.BindIntPath[model.IngestIDType](req, "ingest_id")
	if fieldError != nil {
		response.WriteJSON(
			responseWriter,
			http.StatusBadRequest,
			apierror.NewValidationErrorFromPointer(fieldError),
		)
		return
	}

	params, fieldErrors := bindParams(req.URL.Query())
	if fieldErrors != nil {
		response.WriteJSON(
			responseWriter,
			http.StatusBadRequest,
			apierror.NewValidationError(fieldErrors),
		)
		return
	}

	result, serviceError := hdl.service.GetContext(req.Context(), view.ContextArgs{
		IngestID:       ingestID,
		ExcludeFormats: params.ExcludeFormats,
	})
	if serviceError != nil {
		httperror.WriteServiceError(responseWriter, serviceError)
		return
	}

	if params.IsRich {
		_ = response.WriteJSON(
			responseWriter,
			http.StatusOK,
			mapRichContextResponse(result, hdl.mediaURLBuilder),
		)
	} else {
		_ = response.WriteJSON(
			responseWriter,
			http.StatusOK,
			mapSummaryContextResponse(result),
		)
	}
}
