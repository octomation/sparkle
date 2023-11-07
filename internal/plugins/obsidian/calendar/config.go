package calendar

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/errors"
)

const config = ".obsidian/plugins/calendar/data.json"

var defaults = Config{
	ShouldConfirmBeforeCreate: true,
	WeekStart:                 "locale",
	WordsPerDot:               250,
	ShowWeeklyNote:            false,
	WeeklyNoteFormat:          "gggg-[W]ww",
	WeeklyNoteTemplate:        "",
	WeeklyNoteFolder:          "",
	LocaleOverride:            "system-default",
}

func LoadConfig(fs afero.Fs) (Config, error) {
	cnf := defaults

	f, err := fs.Open(config)
	if os.IsNotExist(err) {
		return cnf, nil
	}
	if err != nil {
		return cnf, errors.X{
			User:   configError,
			System: fmt.Errorf("cannot load config %w", err),
		}
	}
	defer safe.Close(f, unsafe.Ignore)

	if err := json.NewDecoder(f).Decode(&cnf); err != nil {
		return cnf, errors.X{
			User:   configError,
			System: fmt.Errorf("cannot decode config %w", err),
		}
	}

	cnf.enabled = true
	return cnf, nil
}

type Config struct {
	ShouldConfirmBeforeCreate bool   `json:"shouldConfirmBeforeCreate"`
	WeekStart                 string `json:"weekStart"`
	WordsPerDot               int    `json:"wordsPerDot"`
	ShowWeeklyNote            bool   `json:"showWeeklyNote"`
	WeeklyNoteFormat          string `json:"weeklyNoteFormat"`
	WeeklyNoteTemplate        string `json:"weeklyNoteTemplate"`
	WeeklyNoteFolder          string `json:"weeklyNoteFolder"`
	LocaleOverride            string `json:"localeOverride"`

	enabled bool
}

func (Config) Documentation() string {
	return "https://github.com/liamcain/obsidian-calendar-plugin"
}

func (Config) Enabler() string {
	return "Settings > Options > Community plugins > Installed plugins > Calendar"
}

func (cnf Config) IsEnabled() bool {
	return cnf.enabled
}

func (Config) Section() string {
	return "Settings > Community plugins > Calendar"
}

func (Config) ShouldConfirmBeforeCreatePath() string {
	return "> General Settings > Confirm before creating new note"
}

func (Config) WeekStartPath() string {
	return "> General Settings > Start week on"
}

func (Config) WordsPerDotPath() string {
	return "> General Settings > Words per dot"
}

func (Config) ShowWeeklyNotePath() string {
	return "> General Settings > Show week number"
}

// WeeklyNoteFormatPath returns the path to the option
// that allows to change the weekly note filename format.
// See https://momentjs.com/docs/#/displaying/format/ and
// https://github.com/liamcain/obsidian-calendar-plugin#how-do-i-include-unformatted-words-in-my-weekly-note-filenames.
//
// Deprecated: https://github.com/liamcain/obsidian-calendar-plugin?tab=readme-ov-file#weekly-notes-deprecated.
func (Config) WeeklyNoteFormatPath() string {
	return "> Weekly Note Settings > Weekly note format"
}

// Deprecated: https://github.com/liamcain/obsidian-calendar-plugin?tab=readme-ov-file#weekly-notes-deprecated.
func (Config) WeeklyNoteTemplatePath() string {
	return "> Weekly Note Settings > Weekly note template"
}

// Deprecated: https://github.com/liamcain/obsidian-calendar-plugin?tab=readme-ov-file#weekly-notes-deprecated.
func (Config) WeeklyNoteFolderPath() string {
	return "> Weekly Note Settings > Weekly note folder"
}

func (Config) LocaleOverridePath() string {
	return "> Advanced Settings > Override locale"
}
