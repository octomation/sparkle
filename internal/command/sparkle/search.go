package sparkle

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"go.octolab.org/pointer"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"
)

func Search() *cobra.Command {
	const (
		collection = "sparkles"
		ext        = ".md"
		limit      = 3
	)

	var (
		apiKey = os.Getenv("TYPESENSE_API_KEY")
		server = "http://localhost:8108"
	)

	cmd := &cobra.Command{
		Use:  "search [-k api-key] [-s server-address] query",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := typesense.NewClient(
				typesense.WithAPIKey(apiKey),
				typesense.WithServer(server),
			)

			params := &api.SearchCollectionParams{
				Q:                   args[0],
				QueryBy:             "content",
				TypoTokensThreshold: pointer.ToInt(5),
				UseCache:            pointer.ToBool(false),
			}
			result, err := client.Collection(collection).Documents().Search(params)
			if err != nil {
				return err
			}

			if result.Hits == nil || len(*result.Hits) == 0 {
				return errors.New("nothing found")
			}

			for _, item := range *result.Hits {
				doc := *item.Document
				fmt.Println(doc["id"], doc["path"])
				if item.Highlights != nil {
					for _, highlight := range *item.Highlights {
						fmt.Println(
							"\t",
							pointer.ValueOfString(highlight.Field),
							pointer.ValueOfString(highlight.Snippet),
						)
					}
				}
			}
			return nil
		},
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:    "index",
			Args:   cobra.NoArgs,
			Hidden: true,

			RunE: func(cmd *cobra.Command, args []string) error {
				client := typesense.NewClient(
					typesense.WithAPIKey(apiKey),
					typesense.WithServer(server),
				)

				unsafe.DoSilent(client.Collection(collection).Delete())
				schema := &api.CollectionSchema{
					Name: collection,
					Fields: []api.Field{
						{
							Name: "content",
							Type: "string",
						},
						{
							Name: "path",
							Type: "string",
						},
					},
				}
				if _, err := client.Collections().Create(schema); err != nil {
					return err
				}

				fs := afero.NewOsFs()
				root := "."
				files := make([]string, 0, 100)
				err := afero.Walk(fs, root, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					if path == root {
						return nil
					}

					depth := len(strings.Split(
						strings.TrimPrefix(path, root),
						string(filepath.Separator),
					))
					if depth > limit {
						return filepath.SkipDir
					}

					if !info.IsDir() && filepath.Ext(path) == ext {
						files = append(files, path)
					}
					return nil
				})
				if err != nil {
					return err
				}

				type document struct {
					ID      string `json:"id,omitempty"`
					Path    string `json:"path"`
					Content string `json:"content"`
				}

				docs := client.Collection(collection).Documents()
				for _, name := range files {
					file, err := os.Open(name)
					if err != nil {
						return err
					}

					doc, err := io.ReadAll(file)
					if err != nil {
						safe.Close(file, unsafe.Ignore)
						return err
					}
					safe.Close(file, unsafe.Ignore)

					idx, err := docs.Create(document{
						Path:    name,
						Content: string(doc),
					})
					if err != nil {
						return err
					}
					cmd.Println(idx["id"])
				}

				return nil
			},
		},
	)

	flags := cmd.PersistentFlags()
	flags.StringVarP(&apiKey, "api-key", "k", apiKey, "Typesense API key")
	flags.StringVarP(&server, "server", "s", server, "Typesense server address")

	return cmd
}
