package variant

import (
	"strings"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/bind"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
)

func BindImageFormatsQuery(key string, values []string) ([]media.ImageFormat, *apierror.FieldError) {
	rawFormats, fieldError := bind.BindStringSliceQuery(key, values, "+")
	if fieldError != nil {
		return nil, fieldError
	}

	formats := make([]media.ImageFormat, 0, len(rawFormats))
	for _, rawFormat := range rawFormats {
		format, ok := media.ParseImageFormat(strings.ToLower(rawFormat))
		if !ok {
			return nil, &apierror.FieldError{
				Type:    apierror.FieldErrorTypeInvalid,
				Path:    "query." + key,
				Message: "must be a valid image format",
			}
		}

		formats = append(formats, format)
	}

	return formats, nil
}
