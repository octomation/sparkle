package main

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/go-github/v56/github"
	"github.com/google/uuid"
	"go.octolab.org/pointer"
)

func Repositories(
	ctx context.Context,
	cnf Config,
	client *github.Client,
) ([]*github.Repository, error) {
	all := make([]*github.Repository, 0, 10)
	opts := &github.RepositoryListOptions{
		Affiliation: strings.Join(cnf.GitHub.Affiliation, ","),
		ListOptions: github.ListOptions{PerPage: cnf.GitHub.Limit},
	}

	for {
		repos, resp, err := client.Repositories.List(ctx, "", opts)
		if err != nil {
			return nil, err
		}
	filter:
		for _, repo := range repos {
			full := repo.GetFullName()
			if _, present := cnf.GitHub.Include.list[full]; present {
				all = append(all, repo)
				continue
			}

			if _, present := cnf.GitHub.Exclude.list[full]; present {
				continue
			}
			if repo.GetFork() {
				continue
			}
			for _, pattern := range cnf.GitHub.Exclude.patterns {
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

func Dump(
	ctx context.Context,
	cnf Config,
	client *github.Client,
	repos []*github.Repository,
) error {
	collisions := make(map[string]int)
	for _, repo := range repos {
		collisions[repo.GetName()]++
	}

	opts := &github.RepositoryContentGetOptions{}
	for _, repo := range repos {
		owner := repo.GetOwner().GetLogin()
		canonical := cnf.GitHub.canonicals[owner]
		if canonical == "" {
			canonical = owner
		}

		name := repo.GetName()
		full := repo.GetFullName()
		desc := repo.GetDescription()
		home := repo.GetHomepage()
		url := repo.GetHTMLURL()

		ids := []string{strconv.FormatInt(repo.GetID(), 10)}
		// TODO:feature now they are empty, plan the future implementation
		for _, rel := range []*github.Repository{
			repo.GetParent(),
			repo.GetSource(),
			repo.GetTemplateRepository(),
		} {
			if rel != nil {
				ids = append(ids, strconv.FormatInt(rel.GetID(), 10))
			}
		}

		aliases := []string{
			name,
			full,
			strings.TrimPrefix(url, "https://"),
		}
		if owner != canonical {
			aliases = append(aliases, strings.Join([]string{canonical, name}, "/"))
		}
		if extra := cnf.GitHub.aliases[repo.GetSSHURL()]; len(extra) > 0 {
			aliases = append(aliases, extra...)
		}
		tags := cnf.GitHub.Tags
		if repo.GetArchived() {
			tags = append(tags, "archive")
		}

		dir := name
		if collisions[dir] > 1 {
			dir = strings.Join([]string{dir, "at", canonical}, " ")
		}
		dir = filepath.Join(pointer.ValueOfString(dst), dir)

		readme, resp, err := client.Repositories.GetReadme(ctx, owner, name, opts)
		if err != nil {
			if resp.StatusCode != http.StatusNotFound {
				return err
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
			return err
		}

		md := Markdown{
			FrontMatter: FrontMatter{
				InternalID:  uuid.New(),
				ExternalID:  ids,
				Aliases:     aliases,
				Tags:        tags,
				Topics:      repo.Topics,
				Description: desc,
				URL:         url,
				Homepage:    home,
			},
			Content: content,
		}
		if err := md.DumpTo(filepath.Join(dir, cnf.GitHub.Readme)); err != nil {
			return err
		}
	}
	return nil
}
