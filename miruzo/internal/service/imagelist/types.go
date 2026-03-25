package imagelist

import (
	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/samber/mo"
)

type Params[ScalarType model.ImageListCursorScalar] struct {
	// Cursor is the opaque cursor of the last item from the previous page.
	// Use None on the first page.
	Cursor mo.Option[model.ImageListCursorKey[ScalarType]]
	// Limit is the maximum number of items to return.
	Limit uint16
	// ExcludeFormats lists variant formats to exclude.
	// Nil keeps the default format policy.
	ExcludeFormats []media.ImageFormat
}

type Result[ScalarType model.ImageListCursorScalar] struct {
	// Items is the current page of image summaries.
	Items []model.Image
	// Cursor is the key for the next page. None means there is no next page.
	Cursor mo.Option[model.ImageListCursorKey[ScalarType]]
}
