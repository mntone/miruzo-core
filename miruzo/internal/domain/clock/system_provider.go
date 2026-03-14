package clock

import "time"

type systemProvider struct{}

func NewSystemProvider() Provider {
	return systemProvider{}
}

func (systemProvider) Now() time.Time {
	return time.Now().UTC()
}
