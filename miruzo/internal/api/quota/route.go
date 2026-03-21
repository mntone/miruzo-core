package quota

import (
	"net/http"

	m "github.com/mntone/miruzo-core/miruzo/internal/api/middleware"
)

func RegisterRoutes(mux *http.ServeMux, handler handler) {
	mux.HandleFunc("/api/quota",
		m.RequireAcceptJson(
			m.RequireMethodGet(handler.getQuota),
		),
	)
}
