package imagelist

import (
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

type Params[C persist.ImageListCursor] struct {
	Cursor         mo.Option[C]
	Limit          uint16
	ExcludeFormats []string
}

type Result[C persist.ImageListCursor] struct {
	Items  []persist.Image
	Cursor mo.Option[C]
}
