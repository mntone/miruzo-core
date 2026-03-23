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

	mux.HandleFunc("/api/i/{ingest_id}/hall_of_fame",
		cors.New(
			m.CORSAllowMethods{
				http.MethodPost:   m.CORSAllowHeaders{},
				http.MethodDelete: m.CORSAllowHeaders{},
			},
			m.RequireAcceptJson(
				m.RequireMethodsOf(
					[]string{http.MethodPost, http.MethodDelete},
					handler.hallOfFame,
				),
			),
		),
	)
}
