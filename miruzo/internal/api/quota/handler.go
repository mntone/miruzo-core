package quota

import (
	"net/http"

	"github.com/mntone/miruzo-core/miruzo/internal/api/httperror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/response"
	"github.com/mntone/miruzo-core/miruzo/internal/service/user"
)

type handler struct {
	service user.Service
}

func NewHandler(
	srv user.Service,
) *handler {
	return &handler{
		service: srv,
	}
}

func (hdl *handler) getQuota(
	responseWriter http.ResponseWriter,
	req *http.Request,
) {
	result, serviceError := hdl.service.GetQuota(req.Context())
	if serviceError != nil {
		httperror.WriteServiceError(responseWriter, serviceError)
		return
	}

	res, mapError := mapQuota(result)
	if mapError != nil {
		httperror.WriteInternalServerError(responseWriter)
		return
	}

	_ = response.WriteJSON(responseWriter, http.StatusOK, res)
}
