package tact

import (
	"fmt"
	"regexp"
	"time"

	xtime "go.octolab.org/time"
)

const (
	threshold = 12 * time.Hour
	watch     = "15:04"

	_ = iota
	activated
	deactivated
)

var (
	rec  = regexp.MustCompile(`^- (\d{2}:\d{2}) /(?: ((?:\d+[hms])+) /)? (\d{2}:\d{2}) - .*$`)
	end  = regexp.MustCompile(`^\w+ total / \w+ break \d+% / \w+ work \d+%`)
	zero = regexp.MustCompile(`(\D)0[m,s]`)
)

type Logbook struct {
	state  int
	shift  time.Duration
	start  time.Time
	end    time.Time
	breaks time.Duration
}

func (log *Logbook) isActivated() bool {
	return log.state == activated
}

func (log *Logbook) isDeactivated() bool {
	return log.state == deactivated
}

func (log *Logbook) Log(record string) error {
	if log.isDeactivated() || record == "" {
		return nil
	}

	// expected: {record, from, breaks, to} or {exit}
	marks := rec.FindStringSubmatch(record)
	if len(marks) != 4 {
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
	from, err := time.Parse(watch, marks[1])
	if err != nil {
		return err
	}
	from = from.Add(log.shift)
	if log.start.IsZero() {
		log.start = from
	} else if from.Before(log.end) {
		if log.end.Sub(from) < threshold {
			return errTimeTravel
		}
		shift = xtime.Day
		from = from.Add(shift)
	}

	var breaks time.Duration
	if marks[2] != "" {
		breaks, err = time.ParseDuration(marks[2])
		if err != nil {
			return err
		}
	}
	if !log.end.IsZero() {
		breaks += from.Sub(log.end)
	}

	to, err := time.Parse(watch, marks[3])
	if err != nil {
		return err
	}
	to = to.Add(log.shift)
	if from.After(to) {
		if from.Sub(to) < threshold {
			return errTimeTravel
		}
		shift = xtime.Day
		to = to.Add(shift)
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

func (log *Logbook) clean(d time.Duration) string {
	base := d.String()
	iter := zero.ReplaceAllString(base, "$1")
	for iter != base {
		base = iter
		iter = zero.ReplaceAllString(base, "$1")
	}
	return iter
}
