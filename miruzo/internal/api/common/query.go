package common

import "github.com/samber/mo"

type PaginationQuery[Cursor any] struct {
	// Limit is the maximum number of items to return for this request.
	Limit uint16
	// Cursor is an opaque pagination cursor representing the last item returned; nil for the first page.
	Cursor mo.Option[Cursor]
}
