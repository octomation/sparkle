package tact

import (
	"fmt"
	"regexp"
	"time"
)

const (
	watch = "15:04"
)

var (
	rec  = regexp.MustCompile(`^- (\d{2}:\d{2}) /(?: ((?:\d+[hms])+) /)? (\d{2}:\d{2}) - .*$`)
	zero = regexp.MustCompile(`(\D)0[m,s]`)
)

type Logbook struct {
	start  time.Time
	end    time.Time
	breaks time.Duration
}

func (log *Logbook) Log(record string) error {
	// expected: record, from, breaks, to
	if record == "" {
		return nil
	}
	marks := rec.FindStringSubmatch(record)
	if len(marks) != 4 {
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

	from, err := time.Parse(watch, marks[1])
	if err != nil {
		return err
	}
	if log.start.IsZero() {
		log.start = from
	} else if from.Before(log.end) {
		return errTimeTravel
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
	if from.After(to) {
		return errTimeTravel
	}

	log.breaks += breaks
	log.end = to
	return nil
}

func (log *Logbook) Total() time.Duration {
	return log.end.Sub(log.start)
}

func (log *Logbook) Breaks() time.Duration {
	return log.breaks
}

func (log *Logbook) String() string {
	work := log.Total() - log.Breaks()
	rate := int(100 * work / log.Total())
	return fmt.Sprintf(
		"%s total / %s break %d%% / %s work %d%%",
		log.clean(log.Total()),
		log.clean(log.Breaks()), 100-rate,
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
