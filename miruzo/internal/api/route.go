package api

import (
	"net/http"

	httperror "github.com/mntone/miruzo-core/miruzo/internal/api/http/error"
	m "github.com/mntone/miruzo-core/miruzo/internal/api/middleware"
)

func notFound(responseWriter http.ResponseWriter, _ *http.Request) {
	httperror.WriteNotFound(responseWriter)
}

func RegisterNotFoundRoute(mux *http.ServeMux) {
	mux.HandleFunc("/api/", m.RequireAcceptJson(notFound))
}
