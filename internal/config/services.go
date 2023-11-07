package config

import (
	"go.octolab.org/ecosystem/sparkle/internal/plugins/obsidian/calendar"
	diary "go.octolab.org/ecosystem/sparkle/internal/plugins/obsidian/daily-notes"
	"go.octolab.org/ecosystem/sparkle/internal/plugins/obsidian/periodic-notes"
)

type Server struct {
	Name    string  `json:"name"`
	Service Service `json:"service"`
}

type Service struct {
	License string  `json:"license"`
	Plugins Plugins `json:"plugins"`
}

type Plugins struct {
	Obsidian ObsidianPlugins `json:"obsidian"`
}

type ObsidianPlugins struct {
	Calendar      calendar.Config `json:"calendar"`
	DailyNotes    diary.Config    `json:"daily-notes"`
	PeriodicNotes periodic.Config `json:"periodic-notes"`
}
