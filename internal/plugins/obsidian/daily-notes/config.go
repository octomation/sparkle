package diary

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/errors"
)

const config = ".obsidian/daily-notes.json"

var defaults = Config{
	Autorun:  false,
	Folder:   "",
	Format:   "YYYY-MM-DD",
	Template: "",
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

	return cnf, nil
}

type Config struct {
	Autorun  bool   `json:"autorun"`
	Folder   string `json:"folder"`
	Format   string `json:"format"`
	Template string `json:"template"`
}

func (cnf Config) Enabler() string {
	return "Settings > Options > Core plugins > Daily notes"
}

func (cnf Config) Section() string {
	return "Settings > Core plugins > Daily notes"
}

func (cnf Config) AutorunOptionPath() string {
	return "> Open daily note on startup"
}

func (cnf Config) FolderOptionPath() string {
	return "> New file location"
}

// FormatOptionPath returns the path to the option
// that allows to change the date format.
// See https://momentjs.com/docs/#/displaying/format/.
func (cnf Config) FormatOptionPath() string {
	return "> Date format"
}

func (cnf Config) TemplateOptionPath() string {
	return "> Template file location"
}
