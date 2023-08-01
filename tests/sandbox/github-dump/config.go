package main

import (
	"os"
	"regexp"

	"github.com/BurntSushi/toml"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"
)

func Load(name string) (Config, error) {
	cnf := defaults()

	f, err := os.Open(name)
	if err != nil {
		return cnf, err
	}
	defer safe.Close(f, unsafe.Ignore)

	if _, err := toml.NewDecoder(f).Decode(&cnf); err != nil {
		return cnf, err
	}
	return cnf, cnf.Index()
}

type Config struct {
	GitHub struct {
		Affiliation []string `toml:"affiliation"`
		Tags        []string `toml:"tags"`
		Readme      string   `toml:"readme"`
		Limit       int      `toml:"limit"`

		Canonical []struct {
			Name string `toml:"name"`
			Real string `toml:"real"`
		} `toml:"canonical"`
		canonicals map[string]string

		Aliases []struct {
			SSH     string   `toml:"ssh"`
			Aliases []string `toml:"aliases"`
		} `toml:"aliases"`
		aliases map[string][]string

		Include struct {
			List []string `toml:"list"`
			list map[string]struct{}
		} `toml:"include"`

		Exclude struct {
			List []string `toml:"list"`
			list map[string]struct{}

			Patterns []string `toml:"patterns"`
			patterns []*regexp.Regexp
		} `toml:"exclude"`
	} `toml:"github"`
}

func (cnf *Config) Index() error {
	cnf.GitHub.canonicals = make(map[string]string)
	for _, canonical := range cnf.GitHub.Canonical {
		cnf.GitHub.canonicals[canonical.Real] = canonical.Name
	}

	cnf.GitHub.aliases = make(map[string][]string)
	for _, alias := range cnf.GitHub.Aliases {
		cnf.GitHub.aliases[alias.SSH] = alias.Aliases
	}

	cnf.GitHub.Include.list = make(map[string]struct{})
	for _, include := range cnf.GitHub.Include.List {
		cnf.GitHub.Include.list[include] = struct{}{}
	}

	cnf.GitHub.Exclude.list = make(map[string]struct{})
	for _, exclude := range cnf.GitHub.Exclude.List {
		cnf.GitHub.Exclude.list[exclude] = struct{}{}
	}

	size := len(cnf.GitHub.Exclude.Patterns)
	cnf.GitHub.Exclude.patterns = make([]*regexp.Regexp, 0, size)
	for _, pattern := range cnf.GitHub.Exclude.Patterns {
		r, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}
		cnf.GitHub.Exclude.patterns = append(cnf.GitHub.Exclude.patterns, r)
	}
	return nil
}

func defaults() Config {
	var cnf Config
	cnf.GitHub.Affiliation = []string{"owner"}
	cnf.GitHub.Tags = []string{"github", "repository"}
	cnf.GitHub.Readme = "README.md"
	cnf.GitHub.Limit = 100

	// TODO:feature future implementation
	// some names from differed orgs are matched
	collision := "{{ .GetFullName | title }}"
	// sometimes readme are missed
	readme := struct {
		name string
		alt  struct {
			tags    []string
			content string
		}
	}{
		name: "README.md",
		alt: struct {
			tags    []string
			content string
		}{
			tags: []string{"debt", "todo"},
			content: `
# {{ .GetName | emoji }}

{{ .GetDescription }}
`,
		},
	}
	_, _ = collision, readme

	return cnf
}
