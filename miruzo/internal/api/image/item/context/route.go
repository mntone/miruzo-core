package context

import (
	"net/http"

	m "github.com/mntone/miruzo-core/miruzo/internal/api/middleware"
)

func RegisterRoutes(mux *http.ServeMux, handler *handler) {
	mux.HandleFunc("/api/i/{ingest_id}",
		m.RequireAcceptJson(
			m.RequireMethodGet(handler.getContext),
		),
	)
}
