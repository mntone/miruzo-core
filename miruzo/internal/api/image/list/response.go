package list

import (
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type ImageListResponse[C persist.ImageListCursor] struct {
	Items  []ImageListModel `json:"items"`
	Cursor *C               `json:"cursor,omitempty"`
}
