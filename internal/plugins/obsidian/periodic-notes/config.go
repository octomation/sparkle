package periodic

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/errors"
)

const config = ".obsidian/plugins/periodic-notes/data.json"

var defaults = Config{
	ShowGettingStartedBanner:      true,
	HasMigratedDailyNoteSettings:  false,
	HasMigratedWeeklyNoteSettings: false,
	Daily: Period{
		Format:  "YYYY-MM-DD",
		Enabled: false,
	},
	Weekly: Period{
		Format:  "gggg-[W]ww",
		Enabled: false,
	},
	Monthly: Period{
		Format:  "YYYY-MM",
		Enabled: false,
	},
	Quarterly: Period{
		Format:  "YYYY-[Q]Q",
		Enabled: false,
	},
	Yearly: Period{
		Format:  "YYYY",
		Enabled: false,
	},
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

	cnf.Enabled = true
	return cnf, nil
}

type Config struct {
	ShowGettingStartedBanner      bool   `json:"showGettingStartedBanner"`
	HasMigratedDailyNoteSettings  bool   `json:"hasMigratedDailyNoteSettings"`
	HasMigratedWeeklyNoteSettings bool   `json:"hasMigratedWeeklyNoteSettings"`
	Daily                         Period `json:"daily"`
	Weekly                        Period `json:"weekly"`
	Monthly                       Period `json:"monthly"`
	Quarterly                     Period `json:"quarterly"`
	Yearly                        Period `json:"yearly"`

	Enabled bool
}

type Period struct {
	Format   string `json:"format"`
	Folder   string `json:"folder"`
	Template string `json:"template"`
	Enabled  bool   `json:"enabled"`
}

func (Config) Documentation() string {
	return "https://github.com/liamcain/obsidian-periodic-notes"
}

func (Config) Enabler() string {
	return "Settings > Options > Community plugins > Installed plugins > Periodic Notes"
}

func (Config) Section() string {
	return "Settings > Community plugins > Periodic Notes"
}

func (Config) DailyPath() string {
	return "> Daily Notes"
}

func (Config) WeeklyPath() string {
	return "> Weekly Notes"
}

func (Config) MonthlyPath() string {
	return "> Monthly Notes"
}

func (Config) QuarterlyPath() string {
	return "> Quarterly Notes"
}

func (Config) YearlyPath() string {
	return "> Yearly Notes"
}
