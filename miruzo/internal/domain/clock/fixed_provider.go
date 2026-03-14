package clock

import "time"

type fixedProvider struct {
	currentTime time.Time
}

func NewFixedProvider(currentTime time.Time) Provider {
	return fixedProvider{
		currentTime: currentTime,
	}
}

func (clk fixedProvider) Now() time.Time {
	return clk.currentTime
}
