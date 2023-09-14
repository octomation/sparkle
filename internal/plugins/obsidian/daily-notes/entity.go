package diary

import (
	"time"

	"github.com/nleeper/goment"
	xtime "go.octolab.org/time"
)

type Record struct {
	Day    goment.Goment
	Path   string
	Format string
}

func (r Record) Yesterday() *goment.Goment {
	// TODO:debt bad goment design
	//  LinkPrev() and LinkNext() had unexpected side effects
	copied := r.Day
	return copied.Add(-xtime.Day)
}

func (r Record) Tomorrow() *goment.Goment {
	copied := r.Day
	return copied.Add(+xtime.Day)
}

func (r Record) Time() time.Time {
	return r.Day.ToTime()
}
