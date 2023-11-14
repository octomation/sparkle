package tact

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogbook_Log(t *testing.T) {
	tests := map[string]struct {
		logs   []string
		report string
	}{
		"long internal interval": {
			logs: []string{
				"- 12:15 / 13:00 - routine solving / 🥱",
				"- 13:15 / 14:00 - day planning / 🤔",
				"- 14:00 / 15:00 - micro-tasking / 🥱",
				"- 15:00 / 2h15m / 22:15 - focusing on the goal / 😤",
			},
			report: "10h total / 2h30m break 25% / 7h30m work 75%",
		},
		"long working day": {
			logs: []string{
				"- 11:15 / 12:15 - day planning / 🤔",
				"- 12:15 / 13:15 - task solving / 😤",
				"- 13:45 / 45m / 16:30 - reading the book / 😤",
				"- 17:00 / 1h / 21:00 - focusing on the goal / 😬",
				"- 21:00 / 22:00 - solve critical issue / 😬",
				"- 23:00 / 01:15 - write tests / 🫠",
			},
			report: "14h total / 3h45m break 27% / 10h15m work 73%",
		},
		"late start": {
			logs: []string{
				"- 21:00 / 22:00 - solve critical issue / 😬",
				"- 23:00 / 01:15 - write tests / 🫠",
				"- 01:30 / 03:00 - focusing on the goal / 😤",
			},
			report: "6h total / 1h15m break 21% / 4h45m work 79%",
		},
		"two days run": {
			logs: []string{
				"- 11:15 / 12:15 - day planning / 🤔",
				"- 12:15 / 13:15 - task solving / 😤",
				"- 13:45 / 45m / 16:30 - reading the book / 😤",
				"- 17:00 / 1h / 21:00 - focusing on the goal / 😬",
				"- 21:00 / 22:00 - solve critical issue / 😬",
				"- 23:00 / 01:15 - write tests / 🫠",
				"- 01:30 / 08:00 - focusing on the goal / 😤",
				"- 11:15 / 12:15 - day planning / 🤔",
				"- 12:15 / 13:15 - task solving / 😤",
				"- 13:45 / 45m / 16:30 - reading the book / 😤",
				"- 17:00 / 1h / 21:00 - focusing on the goal / 😬",
				"- 21:00 / 22:00 - solve critical issue / 😬",
				"- 23:00 / 01:15 - write tests / 🫠",
			},
			report: "38h total / 11h break 29% / 27h work 71%",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			journal := Logbook{}
			for _, record := range test.logs {
				assert.NoError(t, journal.Log(record))
			}
			assert.Equal(t, test.report, journal.String())
		})
	}
}
