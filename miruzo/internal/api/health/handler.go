package health

import (
	"net/http"

	"github.com/mntone/miruzo-core/miruzo/internal/api/response"
)

type healthResponse struct {
	Version string `json:"version"`
}

type handler struct {
	version string
}

func NewHandler(
	version string,
) handler {
	return handler{
		version: version,
	}
}

func (hdl handler) getHealth(
	responseWriter http.ResponseWriter,
	req *http.Request,
) {
	res := healthResponse{
		Version: hdl.version,
	}
	_ = response.WriteJSON(responseWriter, http.StatusOK, res)
}
