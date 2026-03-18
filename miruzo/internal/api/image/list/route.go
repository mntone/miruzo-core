package list

import (
	"net/http"

	m "github.com/mntone/miruzo-core/miruzo/internal/api/middleware"
)

func RegisterRoutes(mux *http.ServeMux, handler *handler) {
	mux.HandleFunc("/api/i/latest",
		m.RequireAcceptJson(
			m.RequireMethodGet(handler.listLatest),
		),
	)

	mux.HandleFunc("/api/i/chronological",
		m.RequireAcceptJson(
			m.RequireMethodGet(handler.listChronological),
		),
	)

	mux.HandleFunc("/api/i/recently",
		m.RequireAcceptJson(
			m.RequireMethodGet(handler.listRecently),
		),
	)

	mux.HandleFunc("/api/i/first_love",
		m.RequireAcceptJson(
			m.RequireMethodGet(handler.listFirstLove),
		),
	)

	mux.HandleFunc("/api/i/hall_of_fame",
		m.RequireAcceptJson(
			m.RequireMethodGet(handler.listHallOfFame),
		),
	)

	mux.HandleFunc("/api/i/engaged",
		m.RequireAcceptJson(
			m.RequireMethodGet(handler.listEngaged),
		),
	)
}
