package tact

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	xtime "go.octolab.org/ecosystem/sparkle/internal/pkg/time"
)

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

func NewLogbook(format string, threshold time.Duration) *Logbook {
	return &Logbook{format: format, threshold: threshold}
}

type Logbook struct {
	// config
	init      sync.Once
	format    string
	threshold time.Duration

	// log state
	state  int
	shift  time.Duration
	start  time.Time
	end    time.Time
	breaks time.Duration
}

func (log *Logbook) Log(record string) error {
	log.init.Do(func() {
		if log.format == "" {
			log.format = xtime.Kitchen
		}
		if log.threshold == 0 {
			log.threshold = xtime.HalfDay
		}
	})

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

func (log *Logbook) Total() time.Duration {
	return log.end.Sub(log.start)
}

func (log *Logbook) String() string {
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

func (*Logbook) clean(d time.Duration) string {
	base := d.String()
	iter := zero.ReplaceAllString(base, "$1")
	for iter != base {
		base = iter
		iter = zero.ReplaceAllString(base, "$1")
	}
	return iter
}

func (log *Logbook) isActivated() bool {
	return log.state == activated
}

func (log *Logbook) isDeactivated() bool {
	return log.state == deactivated
}
