package reaction

import (
	"net/http"

	m "github.com/mntone/miruzo-core/miruzo/internal/api/middleware"
)

func RegisterRoutes(mux *http.ServeMux, handler *handler) {
	mux.HandleFunc("/api/i/{ingest_id}/love",
		m.RequireAcceptJson(
			m.RequireMethodPost(handler.love),
		),
	)

	mux.HandleFunc("/api/i/{ingest_id}/love/cancel",
		m.RequireAcceptJson(
			m.RequireMethodPost(handler.loveCancel),
		),
	)
}
