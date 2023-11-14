package tact

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogbook_Log(t *testing.T) {
	tests := map[string]struct {
		logs   []string
		total  time.Duration
		breaks time.Duration
		report string
	}{
		"long internal interval": {
			logs: []string{
				"- 12:15 / 13:00 - routine solving / 🥱",
				"- 13:15 / 14:00 - day planning / 🤔",
				"- 14:00 / 15:00 - micro-tasking / 🥱",
				"- 15:00 / 2h15m / 22:15 - focusing on the goal / 😤",
			},
			total:  10 * time.Hour,
			breaks: 2*time.Hour + 30*time.Minute,
			report: "10h total / 2h30m break 25% / 7h30m work 75%",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			journal := Logbook{}
			for _, record := range test.logs {
				assert.NoError(t, journal.Log(record))
			}
			assert.Equal(t, test.total, journal.Total())
			assert.Equal(t, test.breaks, journal.Breaks())
			assert.Equal(t, test.report, journal.String())
		})
	}
}
