package media

import (
	"encoding/json"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestImageFormatJSONRoundTrip(t *testing.T) {
	source := ImageFormatWebP

	encoded, err := json.Marshal(source)
	assert.NilError(t, "json.Marshal(ImageFormatWebP)", err)
	assert.Equal(t, "json marshaled value", string(encoded), `"webp"`)

	var decoded ImageFormat
	err = json.Unmarshal(encoded, &decoded)
	assert.NilError(t, "json.Unmarshal(webp)", err)
	assert.Equal(t, "decoded image format", decoded, ImageFormatWebP)
}

func TestImageFormatUnmarshalJSONReturnsErrorForUnknownFormat(t *testing.T) {
	var decoded ImageFormat
	err := json.Unmarshal([]byte(`"unknown-format"`), &decoded)
	assert.Error(t, "json.Unmarshal(unknown format)", err)
}

func TestImageFormatMarshalJSONReturnsErrorForUnspecifiedValue(t *testing.T) {
	_, err := json.Marshal(ImageFormatUnspecified)
	assert.Error(t, "json.Marshal(ImageFormatUnspecified)", err)
}
