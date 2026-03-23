package middleware

import (
	"net/http"
	"strings"

	httperror "github.com/mntone/miruzo-core/miruzo/internal/api/http/error"
	httpmedia "github.com/mntone/miruzo-core/miruzo/internal/api/http/media"
)

func parseQualityValue(value string) (int32, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}

	intPart := value
	fraction := ""
	if before, after, ok := strings.Cut(value, "."); ok {
		intPart = before
		fraction = after
	}

	if intPart != "0" && intPart != "1" {
		return 0, false
	}
	if len(fraction) > 3 {
		return 0, false
	}

	quality := int32(0)
	digits := int32(0)
	for i := 0; i < len(fraction); i++ {
		ch := fraction[i]
		if ch < '0' || ch > '9' {
			return 0, false
		}
		quality = quality*10 + int32(ch-'0')
		digits++
	}
	for digits < 3 {
		quality *= 10
		digits++
	}

	if intPart == "1" {
		if quality != 0 {
			return 0, false
		}
		return 1000, true
	}

	return quality, true
}

type acceptRange struct {
	httpmedia.MediaType
	quality int32 // 0..1000
}

func parseAcceptRange(segment string) (acceptRange, bool) {
	segment = strings.TrimSpace(segment)
	if segment == "" {
		return acceptRange{}, false
	}

	mediaPart := segment
	paramsPart := ""
	if before, after, ok := strings.Cut(segment, ";"); ok {
		mediaPart = strings.TrimSpace(before)
		paramsPart = after
	}

	parsedMediaType, ok := httpmedia.ParseMediaType(mediaPart)
	if !ok {
		return acceptRange{}, false
	}

	parsedRange := acceptRange{
		MediaType: parsedMediaType,
		quality:   1000,
	}
	if paramsPart == "" {
		return parsedRange, true
	}

	for rawParam := range strings.SplitSeq(paramsPart, ";") {
		param := strings.TrimSpace(rawParam)
		if param == "" {
			continue
		}

		eqIndex := strings.IndexByte(param, '=')
		if eqIndex <= 0 || eqIndex >= len(param)-1 {
			continue
		}

		key := strings.ToLower(strings.TrimSpace(param[:eqIndex]))
		if key != "q" {
			continue
		}

		quality, ok := parseQualityValue(param[eqIndex+1:])
		if !ok {
			return acceptRange{}, false
		}
		parsedRange.quality = quality
		break
	}

	return parsedRange, true
}

func matchSpecificity(accept acceptRange, allowed httpmedia.MediaType) int32 {
	switch {
	case accept.Type == "*" && accept.SubType == "*":
		return 0
	case accept.Type == allowed.Type && accept.SubType == "*":
		return 1
	case accept.Type == allowed.Type && accept.SubType == allowed.SubType:
		return 2
	default:
		return -1
	}
}

type bestMatch struct {
	specificity int32
	quality     int32
}

func isAllowedByAccept(headerValue string, allowedMediaTypes []httpmedia.MediaType) bool {
	if headerValue == "" {
		return true
	}

	bestByAllowedType := make([]bestMatch, len(allowedMediaTypes))
	for i := range allowedMediaTypes {
		bestByAllowedType[i].specificity = -1
	}

	for segment := range strings.SplitSeq(headerValue, ",") {
		accept, ok := parseAcceptRange(segment)
		if !ok {
			continue
		}

		for i, allowed := range allowedMediaTypes {
			specificity := matchSpecificity(accept, allowed)
			if specificity < 0 {
				continue
			}

			if specificity > bestByAllowedType[i].specificity {
				bestByAllowedType[i].specificity = specificity
				bestByAllowedType[i].quality = accept.quality
				continue
			}

			if specificity == bestByAllowedType[i].specificity && accept.quality > bestByAllowedType[i].quality {
				bestByAllowedType[i].quality = accept.quality
			}
		}
	}

	for i := range allowedMediaTypes {
		if bestByAllowedType[i].specificity >= 0 && bestByAllowedType[i].quality > 0 {
			return true
		}
	}
	return false
}

// RequireAcceptAnyOf ensures that the request Accept header allows
// at least one of the specified media types.
//
// Rules:
// - If Accept header is missing or empty: allow (treated as "*/*").
// - If Accept contains "*/*": allow.
// - If Accept contains one of allowedMediaTypes (ignoring parameters): allow.
// - Otherwise: return 406 Not Acceptable.
func RequireAcceptAnyOf(
	allowedMediaTypes []httpmedia.MediaType,
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		acceptHeaderValue := strings.TrimSpace(request.Header.Get("Accept"))
		if isAllowedByAccept(acceptHeaderValue, allowedMediaTypes) {
			next(responseWriter, request)
			return
		}

		httperror.WriteNotAcceptable(responseWriter)
	}
}

// RequireAcceptJson is a shorthand for requiring "application/json".
func RequireAcceptJson(next http.HandlerFunc) http.HandlerFunc {
	return RequireAcceptAnyOf(
		[]httpmedia.MediaType{httpmedia.JSON},
		next,
	)
}
