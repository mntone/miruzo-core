package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/config"
)

func toTotalSeconds(d time.Duration) int64 {
	totalSeconds := int64(d / time.Second)
	return totalSeconds
}

func buildCacheControlHeader(cfg config.StaticFilesConfig) string {
	var parts []string
	if cfg.MaxAge > 0 {
		maxAge := fmt.Sprintf("max-age=%d", toTotalSeconds(cfg.MaxAge))
		parts = append(parts, maxAge)
	}
	if cfg.StaleWhileRevalidate > 0 {
		staleWhileRevalidate := fmt.Sprintf("stale-while-revalidate=%d", toTotalSeconds(cfg.StaleWhileRevalidate))
		parts = append(parts, staleWhileRevalidate)
	}
	if cfg.Immutable {
		parts = append(parts, "immutable")
	}
	if cfg.NoTransform {
		parts = append(parts, "no-transform")
	}
	return strings.Join(parts, ",")
}
