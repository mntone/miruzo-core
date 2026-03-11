package item

import (
	"net/http"
	"net/url"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/bind"
	"github.com/mntone/miruzo-core/miruzo/internal/api/httperror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/response"
	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/service/view"
)

type handler struct {
	service             view.Service
	variantLayersConfig []config.VariantLayerConfig
	mediaURLBuilder     variant.MediaURLBuilder
}

func NewHandler(
	srv view.Service,
	variantLayersConfig []config.VariantLayerConfig,
	mediaURLBuilder variant.MediaURLBuilder,
) *handler {
	return &handler{
		service:             srv,
		variantLayersConfig: variantLayersConfig,
		mediaURLBuilder:     mediaURLBuilder,
	}
}

func parseLevel(qry url.Values) (bool, []apierror.FieldError) {
	if !qry.Has("level") {
		return false, nil
	}

	switch qry.Get("level") {
	case "rich":
		return true, nil
	case "default":
		return false, nil
	}
	return false, []apierror.FieldError{{
		Path:    "query.level",
		Type:    "invalid_value",
		Message: "level must be \"default\" or \"rich\"",
	}}
}

func (hdl *handler) getContext(
	responseWriter http.ResponseWriter,
	req *http.Request,
) {
	ingestID, fieldErrors := bind.ParseIntPath[model.IngestIDType](req, "ingest_id")
	if fieldErrors != nil {
		response.WriteJSON(
			responseWriter,
			http.StatusBadRequest,
			apierror.NewValidationError(fieldErrors),
		)
		return
	}

	rich, fieldErrors := parseLevel(req.URL.Query())
	if fieldErrors != nil {
		response.WriteJSON(
			responseWriter,
			http.StatusBadRequest,
			apierror.NewValidationError(fieldErrors),
		)
		return
	}

	result, serviceError := hdl.service.GetContext(req.Context(), ingestID)
	if serviceError != nil {
		httperror.WriteServiceError(responseWriter, serviceError)
		return
	}

	if rich {
		_ = response.WriteJSON(
			responseWriter,
			http.StatusOK,
			mapRichContextResponse(
				result,
				hdl.variantLayersConfig,
				hdl.mediaURLBuilder,
			),
		)
	} else {
		_ = response.WriteJSON(
			responseWriter,
			http.StatusOK,
			mapSummaryContextResponse(result),
		)
	}
}
