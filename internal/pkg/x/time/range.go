package xtime

import (
	"time"

	xassert "go.octolab.org/ecosystem/sparkle/internal/pkg/x/assert"
)

func Yesterday(in time.Time) time.Time {
	return in.AddDate(0, 0, -1)
}

func Tomorrow(in time.Time) time.Time {
	return in.AddDate(0, 0, +1)
}

func PrevWeek(in time.Time) time.Time {
	return in.AddDate(0, 0, -7)
}

func NextWeek(in time.Time) time.Time {
	return in.AddDate(0, 0, +7)
}

func PrevMonth(in time.Time) time.Time {
	return in.AddDate(0, -1, 0)
}

func NextMonth(in time.Time) time.Time {
	return in.AddDate(0, +1, 0)
}

func PrevQuarter(in time.Time) time.Time {
	return in.AddDate(0, -3, 0)
}

func NextQuarter(in time.Time) time.Time {
	return in.AddDate(0, +3, 0)
}

func PrevYear(in time.Time) time.Time {
	return in.AddDate(-1, 0, 0)
}

func NextYear(in time.Time) time.Time {
	return in.AddDate(+1, 0, 0)
}

func IsLinear(past, future time.Time, threshold time.Duration) (is bool, shift bool) {
	// hard invariant: timestamps within a day
	xassert.True(func() bool {
		delta := future.Sub(past)
		return (delta >= 0 && delta <= Day) ||
			(delta < 0 && delta >= -Day)
	})

	// invariant: past <= future
	if !past.After(future) {
		return true, false
	}

	// invariant: breaks and work time are less than a threshold
	return future.Add(Day).Sub(past) < threshold, true
}
