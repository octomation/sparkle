package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/v56/github"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const limit = 100

type FrontMatter struct {
	InternalID  uuid.UUID `json:"uid" yaml:"uid"`
	ExternalID  int64     `json:"xid" yaml:"xid"`
	Aliases     []string  `json:"aliases" yaml:"aliases"`
	Tags        []string  `json:"tags" yaml:"tags"`
	Topics      []string  `json:"topics" yaml:"topics"`
	Description string    `json:"description" yaml:"description"`
	URL         string    `json:"url" yaml:"url"`
}

type Markdown struct {
	FrontMatter
	Content string
}

func (md Markdown) DumpTo(output io.WriteCloser) error {
	if _, err := fmt.Fprintln(output, "---"); err != nil {
		return err
	}
	enc := yaml.NewEncoder(output)
	enc.SetIndent(2)
	if err := enc.Encode(md.FrontMatter); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(output, "---"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(output, md.Content); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Repository struct {
		Aliases []struct {
			SSH     string   `toml:"ssh"`
			Aliases []string `toml:"aliases"`
		} `toml:"aliases"`
		aliases map[string][]string

		Exclude struct {
			List []string `toml:"list"`
			list map[string]struct{}

			Patterns []string `toml:"patterns"`
			patterns []*regexp.Regexp
		} `toml:"exclude"`
	} `toml:"repository"`
}

func (cnf *Config) Index() error {
	cnf.Repository.aliases = make(map[string][]string)
	for _, alias := range cnf.Repository.Aliases {
		cnf.Repository.aliases[alias.SSH] = alias.Aliases
	}

	cnf.Repository.Exclude.list = make(map[string]struct{})
	for _, exclude := range cnf.Repository.Exclude.List {
		cnf.Repository.Exclude.list[exclude] = struct{}{}
	}

	size := len(cnf.Repository.Exclude.Patterns)
	cnf.Repository.Exclude.patterns = make([]*regexp.Regexp, 0, size)
	for _, pattern := range cnf.Repository.Exclude.Patterns {
		r, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}
		cnf.Repository.Exclude.patterns = append(cnf.Repository.Exclude.patterns, r)
	}
	return nil
}

func main() {
	var cnf Config
	f, err := os.Open("sandbox/github-dump/config.toml")
	if err != nil {
		panic(err)
	}
	if _, err := toml.NewDecoder(f).Decode(&cnf); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := cnf.Index(); err != nil {
		panic(err)
	}

	ctx := context.Background()
	token := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(nil).WithAuthToken(token)

	repos, err := repositories(ctx, cnf, client)
	if err != nil {
		panic(err)
	}

	collisions := make(map[string]int)
	for _, repo := range repos {
		collisions[repo.GetName()]++
	}

	opts := &github.RepositoryContentGetOptions{}
	tags := []string{"github", "repository"}
	for _, repo := range repos {
		owner := repo.GetOwner().GetLogin()
		name := repo.GetName()
		full := repo.GetFullName()
		desc := repo.GetDescription()

		tags := tags // for mutation
		aliases := []string{name, full, strings.TrimPrefix(repo.GetHTMLURL(), "https://")}
		if extra := cnf.Repository.aliases[repo.GetSSHURL()]; len(extra) > 0 {
			aliases = append(aliases, extra...)
		}

		dir := name
		if collisions[dir] > 1 {
			dir = strings.Join([]string{dir, "at", owner}, " ")
		}
		dir = filepath.Join("stream", "github", dir)

		readme, resp, err := client.Repositories.GetReadme(ctx, owner, name, opts)
		if err != nil {
			if resp.StatusCode != http.StatusNotFound {
				panic(err)
			}
			header := name
			tagline := desc
			e, txt := emoji(desc)
			if e != "" {
				header = e + " " + header
				tagline = txt
			}
			tags = append(tags, "debt", "todo")
			draft := fmt.Sprintf("# %s\n\n%s\n", header, tagline)
			readme = &github.RepositoryContent{
				Content: &draft,
			}
		}
		content, err := readme.GetContent()
		if err != nil {
			panic(err)
		}

		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(err)
		}
		f, err := os.Create(filepath.Join(dir, "Repository.md"))
		if err != nil {
			panic(err)
		}
		md := Markdown{
			FrontMatter: FrontMatter{
				InternalID:  uuid.New(),
				ExternalID:  repo.GetID(),
				Aliases:     aliases,
				Tags:        tags,
				Topics:      repo.Topics,
				Description: desc,
				URL:         repo.GetHTMLURL(),
			},
			Content: content,
		}
		if err := md.DumpTo(f); err != nil {
			panic(err)
		}
		if err := f.Close(); err != nil {
			panic(err)
		}
	}
}

func repositories(
	ctx context.Context,
	cnf Config,
	client *github.Client,
) ([]*github.Repository, error) {
	var all []*github.Repository
	opts := &github.RepositoryListOptions{
		Affiliation: "owner,organization_member",
		ListOptions: github.ListOptions{PerPage: limit},
	}
	for {
		repos, resp, err := client.Repositories.List(ctx, "", opts)
		if err != nil {
			return nil, err
		}
	filter:
		for _, repo := range repos {
			full := repo.GetFullName()
			if _, present := cnf.Repository.Exclude.list[full]; present {
				continue
			}
			if repo.GetFork() {
				continue
			}
			for _, pattern := range cnf.Repository.Exclude.patterns {
				if pattern.MatchString(full) {
					continue filter
				}
			}
			all = append(all, repo)
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return all, nil
}

func emoji(in string) (string, string) {
	r, _ := utf8.DecodeRuneInString(in)
	e := string(r)

	// naive implementation
	// research:
	// - https://github.com/spatie/emoji
	// - https://github.com/enescakir/emoji
	emojiRanges := [...]rune{
		0x1F600, 0x1F64F, // Emoticons
		0x1F300, 0x1F5FF, // Misc Symbols and Pictographs
		0x1F680, 0x1F6FF, // Transport and Map
		0x2600, 0x26FF, // Misc symbols
		0x2700, 0x27BF, // Dingbat symbols

		// extended
		0x1FAA0, 0x1FAA8,
	}

	for i := 0; i < len(emojiRanges); i += 2 {
		if r >= emojiRanges[i] && r <= emojiRanges[i+1] {
			return e, strings.TrimSpace(strings.TrimPrefix(in, e))
		}
	}
	return "", in
}
