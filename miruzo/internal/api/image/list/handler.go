package list

import (
	"net/http"

	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/service/imagelist"
)

type handler struct {
	service             imagelist.Service
	variantLayersConfig []config.VariantLayerConfig
	mediaURLBuilder     variant.MediaURLBuilder
}

func NewHandler(
	srv imagelist.Service,
	variantLayersConfig []config.VariantLayerConfig,
	mediaURLBuilder variant.MediaURLBuilder,
) *handler {
	return &handler{
		service:             srv,
		variantLayersConfig: variantLayersConfig,
		mediaURLBuilder:     mediaURLBuilder,
	}
}

func (hdl *handler) listLatest(
	responseWriter http.ResponseWriter,
	req *http.Request,
) {
	hdl.listBase(
		responseWriter,
		req,
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
		hdl.service.ListHallOfFame,
	)
}
