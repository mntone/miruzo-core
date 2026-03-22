package list

import (
	"net/http"

	m "github.com/mntone/miruzo-core/miruzo/internal/api/middleware"
)

func RegisterRoutes(mux *http.ServeMux, cors *m.CORSFactory, handler *handler) {
	mux.HandleFunc("/api/i/latest",
		cors.GET(
			m.RequireAcceptJson(
				m.RequireMethodGet(handler.listLatest),
			),
		),
	)

	mux.HandleFunc("/api/i/chronological",
		cors.GET(
			m.RequireAcceptJson(
				m.RequireMethodGet(handler.listChronological),
			),
		),
	)

	mux.HandleFunc("/api/i/recently",
		cors.GET(
			m.RequireAcceptJson(
				m.RequireMethodGet(handler.listRecently),
			),
		),
	)

	mux.HandleFunc("/api/i/first_love",
		cors.GET(
			m.RequireAcceptJson(
				m.RequireMethodGet(handler.listFirstLove),
			),
		),
	)

	mux.HandleFunc("/api/i/hall_of_fame",
		cors.GET(
			m.RequireAcceptJson(
				m.RequireMethodGet(handler.listHallOfFame),
			),
		),
	)

	mux.HandleFunc("/api/i/engaged",
		cors.GET(
			m.RequireAcceptJson(
				m.RequireMethodGet(handler.listEngaged),
			),
		),
	)
}
