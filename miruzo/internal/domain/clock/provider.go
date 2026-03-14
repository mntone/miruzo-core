package clock

import "time"

type Provider interface {
	Now() time.Time
}
