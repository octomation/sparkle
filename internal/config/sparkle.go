package config

import (
	"go.octolab.org/ecosystem/sparkle/internal/plugins/obsidian/calendar"
	diary "go.octolab.org/ecosystem/sparkle/internal/plugins/obsidian/daily-notes"
	"go.octolab.org/ecosystem/sparkle/internal/plugins/obsidian/periodic-notes"
)

type Sparkle struct {
	Obsidian `json:"obsidian" toml:"obsidian"`
}

type Obsidian struct {
	Plugins struct {
		Calendar      calendar.Config `json:"calendar" toml:"calendar"`
		DailyNotes    diary.Config    `json:"daily-notes" toml:"daily-notes"`
		PeriodicNotes periodic.Config `json:"periodic-notes" toml:"periodic-notes"`
	} `json:"plugins" toml:"plugins"`
}
