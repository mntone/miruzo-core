package api

import (
	"net/http"

	"github.com/mntone/miruzo-core/miruzo/internal/api/httperror"
	m "github.com/mntone/miruzo-core/miruzo/internal/api/middleware"
)

func notFound(responseWriter http.ResponseWriter, _ *http.Request) {
	httperror.WriteNotFound(responseWriter)
}

func RegisterNotFoundRoute(mux *http.ServeMux) {
	mux.HandleFunc("/api/", m.RequireAcceptJson(notFound))
}
