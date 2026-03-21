package media

import (
	"fmt"
	"strings"
)

type MediaType struct {
	Type    string
	SubType string
}

var JSON = MediaType{
	Type:    "application",
	SubType: "json",
}

var ProtocolBuffers = MediaType{
	Type:    "application",
	SubType: "protobuf",
}

func (m MediaType) String() string {
	return fmt.Sprintf("%s/%s", m.Type, m.SubType)
}

func isAlphaNumeric(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9')
}

func isRestrictedNameChar(ch byte) bool {
	if isAlphaNumeric(ch) {
		return true
	}

	switch ch {
	case '!', '#', '$', '&', '-', '^', '_', '.', '+':
		return true
	default:
		return false
	}
}

func isRestrictedName(value string) bool {
	if len(value) == 0 {
		return false
	}
	if !isAlphaNumeric(value[0]) {
		return false
	}

	for i := 1; i < len(value); i++ {
		if !isRestrictedNameChar(value[i]) {
			return false
		}
	}
	return true
}

func ParseMediaType(value string) (MediaType, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return MediaType{}, false
	}

	if semicolonIndex := strings.IndexByte(value, ';'); semicolonIndex >= 0 {
		value = value[:semicolonIndex]
	}
	if value == "" {
		return MediaType{}, false
	}

	mediaType, subType, ok := strings.Cut(value, "/")
	if !ok || mediaType == "" || subType == "" {
		return MediaType{}, false
	}

	if subType == "*" {
		if mediaType == "*" {
			return MediaType{Type: "*", SubType: "*"}, true
		}
		if !isRestrictedName(mediaType) {
			return MediaType{}, false
		}

		return MediaType{
			Type:    strings.ToLower(mediaType),
			SubType: "*",
		}, true
	}
	if mediaType == "*" {
		return MediaType{}, false
	}
	if !isRestrictedName(mediaType) || !isRestrictedName(subType) {
		return MediaType{}, false
	}

	return MediaType{
		Type:    strings.ToLower(mediaType),
		SubType: strings.ToLower(subType),
	}, true
}

func ParseMediaTypes(mediaTypeStrings []string) []MediaType {
	if mediaTypeStrings == nil {
		return nil
	}

	parsedMediaTypes := make([]MediaType, 0, len(mediaTypeStrings))
	for _, mediaType := range mediaTypeStrings {
		parsedMediaType, ok := ParseMediaType(mediaType)
		if ok {
			parsedMediaTypes = append(parsedMediaTypes, parsedMediaType)
		}
	}
	return parsedMediaTypes
}
