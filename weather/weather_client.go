package weather

import "time"

type WeatherClient interface {
	GetWarnings() []Warning
}

type Warning struct {
	Title, Text string
	Begin, End  time.Time
	IsActive    bool
}
