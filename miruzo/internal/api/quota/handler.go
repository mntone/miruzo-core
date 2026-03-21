package quota

import (
	"net/http"

	httperror "github.com/mntone/miruzo-core/miruzo/internal/api/http/error"
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

	_ = response.WriteJSON(responseWriter, http.StatusOK, mapQuota(result))
}
