//go:build dev

package app

import "net/http"

func addHeaders(cacheControl string, noSniff bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		h := responseWriter.Header()
		h.Set("Cache-Control", cacheControl)
		if noSniff {
			h.Set("X-Content-Type-Options", "nosniff")
		}
		next.ServeHTTP(responseWriter, request)
	})
}

func mountStatic(
	mux *http.ServeMux,
	rootDir string,
	basePath string,
	cacheControl string,
	noSniff bool,
) {
	fileServer := http.FileServer(http.Dir(rootDir))
	mux.Handle(
		basePath,
		http.StripPrefix(
			basePath,
			addHeaders(cacheControl, noSniff, fileServer),
		),
	)
}
