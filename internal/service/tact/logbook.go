package tact

import (
	"fmt"
	"regexp"
	"time"

	xtime "go.octolab.org/ecosystem/sparkle/internal/pkg/time"
)

type Logbook interface {
	fmt.Stringer

	Log(string) error
	Total() time.Duration
}

func NewBulletJournal() Logbook {
	panic("not implemented yet, but see https://bulletjournal.com/")
}

func NewLinearJournal(opts ...LinearOption) Logbook {
	journal := &linear{
		format:    xtime.Kitchen,
		threshold: xtime.HalfDay,
	}
	for _, option := range opts {
		option(journal)
	}
	return journal
}

type LinearOption func(*linear)

func WithLinearFormat(format string) LinearOption {
	return func(log *linear) {
		log.format = format
	}
}

func WithLinearThreshold(threshold time.Duration) LinearOption {
	return func(log *linear) {
		log.threshold = threshold
	}
}

const (
	_ = iota
	activated
	deactivated
)

var (
	rec  = regexp.MustCompile(`^- (\d{2}:\d{2}) /(?: ((?:\d+[hms])+) /)? (\d{2}:\d{2}) - .*$`)
	end  = regexp.MustCompile(`^\w+ total / \w+ break \d+% / \w+ work \d+%`)
	zero = regexp.MustCompile(`(\D)0[m,s]`) // 12H
)

type linear struct {
	// config
	format    string
	threshold time.Duration

	// log state
	state  int
	shift  time.Duration
	start  time.Time
	end    time.Time
	breaks time.Duration
}

func (log *linear) Log(record string) error {
	if log.isDeactivated() || record == "" {
		return nil
	}

	// expected: {record, from, breaks, to} or {exit}
	markers := rec.FindStringSubmatch(record)
	if len(markers) != 4 {
		if !log.isActivated() {
			return nil
		}
		if !end.MatchString(record) {
			return nil
		}

		log.state = deactivated
		if expected := log.String(); record != expected {
			return fmt.Errorf(
				"error: %w\nexpected: %s\nobtained: %s",
				errInvalidRecord,
				expected,
				record,
			)
		}
		return nil
	}
	log.state = activated

	var shift time.Duration
	from, err := time.Parse(log.format, markers[1])
	if err != nil {
		return err
	}
	from = from.Add(log.shift)
	if log.start.IsZero() {
		log.start = from
	} else {
		is, shifted := xtime.IsLinear(log.end, from, log.threshold)
		if !is {
			return errTimeTravel
		}
		if shifted {
			shift = xtime.Day
			from = from.Add(shift)
		}
	}

	var breaks time.Duration
	if markers[2] != "" {
		breaks, err = time.ParseDuration(markers[2])
		if err != nil {
			return err
		}
	}

	to, err := time.Parse(log.format, markers[3])
	if err != nil {
		return err
	}
	to = to.Add(log.shift)
	// invariant: work time < threshold, work time = duration(from, to) - breaks
	is, shifted := xtime.IsLinear(from, to, log.threshold+breaks)
	if !is {
		return errTimeTravel
	}
	if shifted {
		shift = xtime.Day
		to = to.Add(shift)
	}

	if !log.end.IsZero() {
		breaks += from.Sub(log.end)
	}
	log.breaks += breaks
	log.shift += shift
	log.end = to
	return nil
}

func (log *linear) Total() time.Duration {
	return log.end.Sub(log.start)
}

func (log *linear) String() string {
	total := log.Total()
	if total == 0 {
		return ""
	}

	work := total - log.breaks
	rate := int(100 * work / total)
	return fmt.Sprintf(
		"%s total / %s break %d%% / %s work %d%%",
		log.clean(total),
		log.clean(log.breaks), 100-rate,
		log.clean(work), rate,
	)
}

func (*linear) clean(d time.Duration) string {
	base := d.String()
	iter := zero.ReplaceAllString(base, "$1")
	for iter != base {
		base = iter
		iter = zero.ReplaceAllString(base, "$1")
	}
	return iter
}

func (log *linear) isActivated() bool {
	return log.state == activated
}

func (log *linear) isDeactivated() bool {
	return log.state == deactivated
}
