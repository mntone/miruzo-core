package variant

import (
	"strings"

	"github.com/mntone/miruzo-core/miruzo/internal/config"
)

type MediaURLBuilder struct {
	basePath string
}

func normalizeBasePath(basePath string) string {
	normalized := strings.TrimSpace(basePath)
	if normalized == "" {
		return "/"
	}
	if !strings.HasPrefix(normalized, "/") {
		normalized = "/" + normalized
	}
	if !strings.HasSuffix(normalized, "/") {
		normalized += "/"
	}
	return normalized
}

func NewMediaURLBuilder(conf config.MediaPublicConfig) MediaURLBuilder {
	return MediaURLBuilder{
		basePath: normalizeBasePath(conf.BasePath),
	}
}

func (builder MediaURLBuilder) Build(relativePath string) string {
	// return builder.basePath + strings.TrimLeft(relativePath, "/")
	return builder.basePath + relativePath
}
