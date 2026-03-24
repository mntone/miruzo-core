//go:build !dev

package app

import "net/http"

func mountStatic(
	_ *http.ServeMux,
	_ string,
	_ string,
	_ string,
	_ bool,
) {
}
