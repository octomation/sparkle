package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"
)

const (
	address = "http://localhost:8108"
	token   = "612fc8ac-ac81-4b79-8651-db55b2783698"
)

func main() {
	client := typesense.NewClient(
		typesense.WithServer(address),
		typesense.WithAPIKey(token),
	)

	collection := "wisdom"
	migrate(client, collection)
	index(client, collection, "stream/wisdom/tg-*.md")
}

func migrate(client *typesense.Client, collection string) {
	if _, err := client.Collection(collection).Delete(); err != nil {
		panic(err)
	}

	schema := &api.CollectionSchema{
		Name: collection,
		Fields: []api.Field{
			{
				Name: "content",
				Type: "string",
			},
			{
				Name: "url",
				Type: "string",
			},
		},
	}
	if _, err := client.Collections().Create(schema); err != nil {
		panic(err)
	}
}

func index(client *typesense.Client, collection, glob string) {
	files, err := filepath.Glob(glob)
	if err != nil {
		panic(err)
	}
	if len(files) == 0 {
		panic(errors.New("there are no files"))
	}

	type document struct {
		ID      string `json:"id"`
		URL     string `json:"url"`
		Content string `json:"content"`
	}

	docs := client.Collection(collection).Documents()
	for _, name := range files {
		file, err := os.Open(name)
		if err != nil {
			panic(err)
		}

		md, err := pageparser.ParseFrontMatterAndContent(file)
		safe.Close(file, unsafe.Ignore)
		if err != nil {
			panic(err)
		}

		doc := document{
			ID:      md.FrontMatter["uid"].(string),
			URL:     md.FrontMatter["url"].(string),
			Content: string(md.Content),
		}
		if _, err := docs.Upsert(doc); err != nil {
			panic(err)
		}
	}
}
