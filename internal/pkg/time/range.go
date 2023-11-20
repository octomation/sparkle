package xtime

import (
	"time"

	xassert "go.octolab.org/ecosystem/sparkle/internal/pkg/assert"
)

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
