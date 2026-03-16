package timezone

import "github.com/samber/mo"

type TimezoneResolver interface {
	GetLocation() mo.Option[string]
}
