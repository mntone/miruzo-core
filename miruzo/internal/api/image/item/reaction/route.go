package reaction

import (
	"net/http"

	m "github.com/mntone/miruzo-core/miruzo/internal/api/middleware"
)

func RegisterRoutes(mux *http.ServeMux, cors *m.CORSFactory, handler handler) {
	mux.HandleFunc("/api/i/{ingest_id}/love",
		cors.POST(
			m.RequireAcceptJson(
				m.RequireMethodPost(handler.love),
			),
		),
	)

	mux.HandleFunc("/api/i/{ingest_id}/love/cancel",
		cors.POST(
			m.RequireAcceptJson(
				m.RequireMethodPost(handler.loveCancel),
			),
		),
	)
}
