package sparkle

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/typesense/typesense-go/typesense"
	"go.octolab.org/pointer"

	"go.octolab.org/ecosystem/sparkle/internal/service/search"
	"go.octolab.org/ecosystem/sparkle/internal/service/search/schema"
	"go.octolab.org/ecosystem/sparkle/internal/service/vault"
)

func Search() *cobra.Command {
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
			service := search.New(client)

			docs, err := service.Find(args[0])
			if err != nil {
				return err
			}
			for _, doc := range docs {
				fmt.Println(doc.ID, ":\t", doc.Path, time.Unix(doc.UpdatedAt, 0).Format(time.RFC3339))
				for _, highlight := range doc.Highlights {
					fmt.Println(
						"\t",
						pointer.ValueOfString(highlight.Field),
						pointer.ValueOfString(highlight.Snippet),
					)
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
				service := search.New(client)

				if err := service.Migrate(true); err != nil {
					return err
				}

				notes, err := vault.New(afero.NewOsFs()).Notes(".")
				if err != nil {
					return err
				}

				docs := make([]schema.Sparkle, 0, len(notes))
				for _, note := range notes {
					docs = append(docs, schema.Sparkle{
						ID:         note.ID(),
						Path:       note.Path,
						Properties: note.Properties,
						Content:    note.Content(),
						CreatedAt:  note.CreatedAt.Unix(),
						UpdatedAt:  note.UpdatedAt.Unix(),
					})
				}
				return service.Index(docs...)
			},
		},
	)

	flags := cmd.PersistentFlags()
	flags.StringVarP(&apiKey, "api-key", "k", apiKey, "Typesense API key")
	flags.StringVarP(&server, "server", "s", server, "Typesense server address")

	return cmd
}
