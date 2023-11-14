package command

import (
	"github.com/spf13/cobra"

	"go.octolab.org/ecosystem/sparkle/internal/command/sparkle"
)

// New returns the new root command.
func New() *cobra.Command {
	command := cobra.Command{
		Use:   "sparkle",
		Short: "✨ Sparkle service.",
		Long:  "The personal development framework and Personal Knowledge Management platform.",

		Args: cobra.NoArgs,

		SilenceErrors: false,
		SilenceUsage:  true,
	}

	command.AddCommand(
		sparkle.Search(),
		sparkle.Stream(),
		sparkle.Tact(),
		NewServer(),
		NewClient(),
	)

	return &command
}
