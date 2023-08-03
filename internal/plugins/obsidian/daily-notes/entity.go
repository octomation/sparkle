package diary

import (
	"time"

	"github.com/nleeper/goment"
)

type Record struct {
	Day    goment.Goment
	Path   string
	Format string
}

func (e Record) Time() time.Time {
	return e.Day.ToTime()
}
