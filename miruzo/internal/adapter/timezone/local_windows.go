// Update windows-to-IANA mapping with:
//   go run ./tools/gen_windows_to_iana/main.go

package timezone

import (
	"github.com/samber/mo"
	"golang.org/x/sys/windows"
)

func getLocalTimezone() string {
	var timezoneInfo DynamicTimeZoneInformation
	_, err := GetDynamicTimeZoneInformation(&timezoneInfo)
	if err != nil {
		return ""
	}

	keyName := windows.UTF16ToString(timezoneInfo.TimeZoneKeyName[:])
	ianaName, ok := windowsToIana[keyName]
	if !ok {
		return ""
	}

	return ianaName
}

type WindowsLocalTimezoneResolver struct {
	location mo.Option[string]
	loaded   bool
}

func NewLocalTimezoneResolver() TimezoneResolver {
	return &WindowsLocalTimezoneResolver{}
}

func (resolv *WindowsLocalTimezoneResolver) GetLocation() mo.Option[string] {
	if !resolv.loaded {
		resolv.loaded = true
		if tz := getLocalTimezone(); tz != "" {
			resolv.location = mo.Some(tz)
		}
	}

	return resolv.location
}
