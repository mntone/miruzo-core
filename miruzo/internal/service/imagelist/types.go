package imagelist

import (
	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

type Params[C persist.ImageListCursor] struct {
	// Cursor is the opaque cursor of the last item from the previous page.
	// Use None on the first page.
	Cursor mo.Option[C]
	// Limit is the maximum number of items to return.
	Limit uint16
	// ExcludeFormats lists variant formats to exclude.
	// Nil keeps the default format policy.
	ExcludeFormats []media.ImageFormat
}

type Result[C persist.ImageListCursor] struct {
	Items  []model.Image
	Cursor mo.Option[C]
}
