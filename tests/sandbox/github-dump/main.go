package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/go-github/v56/github"
	"go.octolab.org/pointer"
)

var (
	config  = flag.String("c", "config.toml", "path to config file")
	dst     = flag.String("d", ".", "path to destination directory")
	replace = flag.Bool("r", false, "overwrite existing files")
)

func main() {
	flag.Parse()
	ctx := context.Background()
	token := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(nil).WithAuthToken(token)

	cnf, err := Load(pointer.ValueOfString(config))
	if err != nil {
		panic(err)
	}

	repos, err := Repositories(ctx, cnf, client)
	if err != nil {
		panic(err)
	}

	if err := Dump(ctx, cnf, client, repos); err != nil {
		panic(err)
	}
}
