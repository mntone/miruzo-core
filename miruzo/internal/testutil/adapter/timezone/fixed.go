package timezone

import (
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/timezone"
	"github.com/samber/mo"
)

type fixedTimezoneResolver struct {
	location mo.Option[string]
}

func NewEmptyTimezoneResolver() timezone.TimezoneResolver {
	return fixedTimezoneResolver{}
}

func NewFixedTimezoneResolverWithLocation(location string) timezone.TimezoneResolver {
	return fixedTimezoneResolver{
		location: mo.Some(location),
	}
}

func (resolv fixedTimezoneResolver) GetLocation() mo.Option[string] {
	return resolv.location
}
