package health

import (
	"net/http"

	m "github.com/mntone/miruzo-core/miruzo/internal/api/middleware"
)

func RegisterRoutes(mux *http.ServeMux, handler handler) {
	mux.HandleFunc("/api/health",
		m.RequireAcceptJson(
			m.RequireMethodGet(handler.getHealth),
		),
	)
}
