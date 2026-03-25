package variant

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestBindImageFormatsQueryReturnsValidValues(t *testing.T) {
	got, err := BindImageFormatsQuery("exclude_formats", []string{"jpeg+webp"})
	assert.Nil(t, "BindImageFormatsQuery() error", err)
	assert.LenIs(t, "BindImageFormatsQuery()", got, 2)
	assert.Equal(t, "BindImageFormatsQuery()[0]", got[0], media.ImageFormatJPEG)
	assert.Equal(t, "BindImageFormatsQuery()[1]", got[1], media.ImageFormatWebP)
}

func TestBindImageFormatsQueryAcceptsMixedCaseValues(t *testing.T) {
	got, err := BindImageFormatsQuery("exclude_formats", []string{"JpEg+WEBP"})
	assert.Nil(t, "BindImageFormatsQuery() error", err)
	assert.LenIs(t, "BindImageFormatsQuery()", got, 2)
	assert.Equal(t, "BindImageFormatsQuery()[0]", got[0], media.ImageFormatJPEG)
	assert.Equal(t, "BindImageFormatsQuery()[1]", got[1], media.ImageFormatWebP)
}

func TestBindImageFormatsQueryReturnsErrorForUnknownFormat(t *testing.T) {
	_, err := BindImageFormatsQuery("exclude_formats", []string{"webp+unknown"})
	assert.NotNil(t, "BindImageFormatsQuery() error", err)
	assert.Equal(t, "BindImageFormatsQuery().Type", err.Type, "invalid")
	assert.Equal(t, "BindImageFormatsQuery().Path", err.Path, "query.exclude_formats")
	assert.Equal(t, "BindImageFormatsQuery().Message", err.Message, "must be a valid image format")
}

func TestBindImageFormatsQueryReturnsErrorForDuplicateQuery(t *testing.T) {
	_, err := BindImageFormatsQuery("exclude_formats", []string{"jpeg", "webp"})
	assert.NotNil(t, "BindImageFormatsQuery() error", err)
	assert.Equal(t, "BindImageFormatsQuery().Type", err.Type, "duplicate")
	assert.Equal(t, "BindImageFormatsQuery().Path", err.Path, "query.exclude_formats")
	assert.Equal(
		t,
		"BindImageFormatsQuery().Message",
		err.Message,
		"must not be specified multiple times",
	)
}
