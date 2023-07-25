package command

import (
	"net/http"

	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"go.octolab.org/ecosystem/sparkle/internal/api/service/v1/servicev1connect"
)

// NewServer returns the new server command.
func NewServer() *cobra.Command {
	command := &cobra.Command{
		Use:   "sparkle",
		Short: "sparkle server",
		Long:  "✨ Sparkle service provider.",

		Args: cobra.NoArgs,

		SilenceErrors: false,
		SilenceUsage:  true,
	}

	command.AddCommand(&cobra.Command{
		Use: "run",

		RunE: func(*cobra.Command, []string) error {
			mux := http.NewServeMux()

			path, handler := servicev1connect.NewServiceHandler(new(service))
			mux.Handle(path, handler)

			return http.ListenAndServe(addr, h2c.NewHandler(mux, new(http2.Server)))
		},
	})

	return command
}
