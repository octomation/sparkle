package tact

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogbook_Log(t *testing.T) {
	tests := map[string]struct {
		logs   []string
		report string
		desc   string
	}{
		"threshold challenge": {
			logs: []string{
				"- 09:30 / 10:00 - day planning and reflection / 🤔",
				"- 10:00 / 5h / 00:00 - focused work on tasks / 🫠",
			},
			report: "14h30m total / 5h break 35% / 9h30m work 65%",
			desc: `
				There is a threshold to prevent incorrect inputs, e.g., from the past.
				An example:
					- 09:30 / 10:00 - day planning / 🤔
					- 09:50 / 12:00 - hard work / 😤
				A primitive solution for checking linearity has a disadvantage:
					- 23:00 / 01:00 - hard work / 😤
				From "23:00" >> To "01:00" because they are parsed for the same day.
				To handle this case, we must define a work time threshold.
			`,
		},
		"long breaks between actions": {
			logs: []string{
				"- 09:15 / 10:00 - day planning / 🤔",
				"- 13:00 / 15:00 - routine solving / 🥱",
				"- 16:00 / 19:15 - goal achieving / 😤",
			},
			report: "10h total / 4h break 40% / 6h work 60%",
		},
		"long breaks inside actions": {
			logs: []string{
				"- 09:15 / 10:00 - day planning / 🤔",
				"- 11:00 / 2h / 15:00 - routine solving / 🥱",
				"- 16:00 / 1h / 19:15 - goal achieving / 😤",
			},
			report: "10h total / 5h break 50% / 5h work 50%",
		},
		"long working day": {
			logs: []string{
				"- 11:15 / 12:15 - day planning / 🤔",
				"- 12:15 / 13:15 - task solving / 😤",
				"- 13:45 / 45m / 16:30 - reading the book / 😤",
				"- 17:00 / 1h / 21:00 - focusing on the goal / 😬",
				"- 21:00 / 22:00 - solve critical issue / 😬",
				"- 23:15 / 01:15 - write tests / 🫠",
			},
			report: "14h total / 4h break 29% / 10h work 71%",
		},
		"late start": {
			logs: []string{
				"- 21:00 / 22:00 - solve critical issue / 😬",
				"- 23:00 / 01:15 - write tests / 🫠",
				"- 01:30 / 45m / 07:00 - focusing on the goal / 😤",
			},
			report: "10h total / 2h break 20% / 8h work 80%",
		},
		"two days run": {
			logs: []string{
				"- 11:15 / 12:15 - day planning / 🤔",
				"- 12:15 / 13:15 - task solving / 😤",
				"- 13:45 / 45m / 16:30 - reading the book / 😤",
				"- 17:00 / 1h / 21:00 - focusing on the goal / 😬",
				"- 21:00 / 22:00 - solve critical issue / 😬",
				"- 23:15 / 01:15 - write tests / 🫠",
				"- 01:30 / 45m / 07:00 - focusing on the goal / 😤",
				"- 11:15 / 12:15 - day planning / 🤔",
				"- 12:15 / 13:15 - task solving / 😤",
				"- 13:45 / 45m / 16:30 - reading the book / 😤",
				"- 17:00 / 1h / 21:00 - focusing on the goal / 😬",
				"- 21:00 / 22:00 - solve critical issue / 😬",
				"- 23:15 / 01:15 - write tests / 🫠",
			},
			report: "38h total / 13h15m break 35% / 24h45m work 65%",
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
