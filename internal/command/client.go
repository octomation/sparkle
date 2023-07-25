package command

import (
	"encoding/json"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	api "go.octolab.org/ecosystem/sparkle/internal/api/service/v1"
	"go.octolab.org/ecosystem/sparkle/internal/api/service/v1/servicev1connect"
)

// NewClient returns the new client command.
func NewClient() *cobra.Command {
	command := &cobra.Command{
		Use:   "sparkle",
		Short: "sparkle client",
		Long:  "✨ Sparkle service controller.",

		Args: cobra.NoArgs,

		SilenceErrors: false,
		SilenceUsage:  true,
	}

	command.AddCommand(&cobra.Command{
		Use: "whoami",

		RunE: func(cmd *cobra.Command, _ []string) error {
			client := servicev1connect.NewServiceClient(
				http.DefaultClient,
				fmt.Sprintf("http://%s", addr),
			)

			resp, err := client.WhoAmI(cmd.Context(), &connect.Request[api.PingRequest]{
				Msg: &api.PingRequest{Env: true},
			})
			if err != nil {
				return err
			}

			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(resp.Msg)
		},
	})

	return command
}
