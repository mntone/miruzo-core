package quota

import (
	"net/http"

	m "github.com/mntone/miruzo-core/miruzo/internal/api/middleware"
)

func RegisterRoutes(mux *http.ServeMux, cors *m.CORSFactory, handler handler) {
	mux.HandleFunc("/api/quota",
		cors.GET(
			m.RequireAcceptJson(
				m.RequireMethodGet(handler.getQuota),
			),
		),
	)
}
