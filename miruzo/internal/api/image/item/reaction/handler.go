package reaction

import (
	"net/http"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/bind"
	"github.com/mntone/miruzo-core/miruzo/internal/api/httperror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/response"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/service/reaction"
)

type handler struct {
	service reaction.Service
}

func NewHandler(
	srv reaction.Service,
) *handler {
	return &handler{
		service: srv,
	}
}

func (hdl *handler) love(
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

	result, serviceError := hdl.service.Love(req.Context(), ingestID)
	if serviceError != nil {
		httperror.WriteServiceError(responseWriter, serviceError)
		return
	}

	res, mapError := mapLoveResponse(result)
	if mapError != nil {
		httperror.WriteInternalServerError(responseWriter)
		return
	}

	_ = response.WriteJSON(responseWriter, http.StatusOK, res)
}

func (hdl *handler) loveCancel(
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

	result, serviceError := hdl.service.LoveCancel(req.Context(), ingestID)
	if serviceError != nil {
		httperror.WriteServiceError(responseWriter, serviceError)
		return
	}

	res, mapError := mapLoveResponse(result)
	if mapError != nil {
		httperror.WriteInternalServerError(responseWriter)
		return
	}

	_ = response.WriteJSON(responseWriter, http.StatusOK, res)
}
