package list

import (
	"net/http"

	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/service/imagelist"
)

type handler struct {
	service           imagelist.Service
	variantLayersSpec media.VariantLayersSpec
	mediaURLBuilder   variant.MediaURLBuilder
}

func NewHandler(
	srv imagelist.Service,
	variantLayersSpec media.VariantLayersSpec,
	mediaURLBuilder variant.MediaURLBuilder,
) *handler {
	return &handler{
		service:           srv,
		variantLayersSpec: variantLayersSpec,
		mediaURLBuilder:   mediaURLBuilder,
	}
}

func (hdl *handler) listLatest(
	responseWriter http.ResponseWriter,
	req *http.Request,
) {
	hdl.listBase(
		responseWriter,
		req,
		imageListCursorModeLatest,
		hdl.service.ListLatest,
	)
}

func (hdl *handler) listChronological(
	responseWriter http.ResponseWriter,
	req *http.Request,
) {
	hdl.listBase(
		responseWriter,
		req,
		imageListCursorModeChronological,
		hdl.service.ListChronological,
	)
}

func (hdl *handler) listRecently(
	responseWriter http.ResponseWriter,
	req *http.Request,
) {
	hdl.listBase(
		responseWriter,
		req,
		imageListCursorModeRecently,
		hdl.service.ListRecently,
	)
}

func (hdl *handler) listFirstLove(
	responseWriter http.ResponseWriter,
	req *http.Request,
) {
	hdl.listBase(
		responseWriter,
		req,
		imageListCursorModeFirstLove,
		hdl.service.ListFirstLove,
	)
}

func (hdl *handler) listHallOfFame(
	responseWriter http.ResponseWriter,
	req *http.Request,
) {
	hdl.listBase(
		responseWriter,
		req,
		imageListCursorModeHallOfFame,
		hdl.service.ListHallOfFame,
	)
}
