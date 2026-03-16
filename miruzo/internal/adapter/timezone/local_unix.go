//go:build unix

package timezone

import (
	"os"
	"strings"

	"github.com/samber/mo"
)

const zoneinfoPrefix = "/usr/share/zoneinfo/"

func getLocalTimezone() string {
	path, err := os.Readlink("/etc/localtime")
	if err != nil {
		return ""
	}

	if !strings.HasPrefix(path, zoneinfoPrefix) {
		return ""
	}

	return strings.TrimPrefix(path, zoneinfoPrefix)
}

type UnixLocalTimezoneResolver struct {
	location mo.Option[string]
	loaded   bool
}

func NewLocalTimezoneResolver() TimezoneResolver {
	return &UnixLocalTimezoneResolver{}
}

func (resolv *UnixLocalTimezoneResolver) GetLocation() mo.Option[string] {
	if !resolv.loaded {
		resolv.loaded = true
		if tz := getLocalTimezone(); tz != "" {
			resolv.location = mo.Some(tz)
		}
	}

	return resolv.location
}
