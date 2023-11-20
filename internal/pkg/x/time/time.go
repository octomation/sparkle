package xtime

import "time"

const (
	Kitchen = "15:04" // mix of the time.Kitchen and time.TimeOnly

	HalfDay = 12 * time.Hour
	Day     = 24 * time.Hour
	Week    = 7 * Day
)
