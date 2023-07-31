package sparkle

import (
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	xtime "go.octolab.org/time"

	diary "go.octolab.org/ecosystem/sparkle/internal/plugins/obsidian/daily-notes"
)

func Stream() *cobra.Command {
	cmd := &cobra.Command{
		Use: "stream",
	}

	cmd.AddCommand(
		Diary(),
	)

	return cmd
}

func Diary() *cobra.Command {
	cmd := &cobra.Command{
		Use: "diary",
	}

	var (
		since   = time.Now().Format(time.DateOnly)
		until   = time.Now().Add(xtime.Week).Format(time.DateOnly)
		rewrite = false
	)
	makeCmd := &cobra.Command{
		Use:  "make",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fs := afero.NewOsFs()
			config, err := diary.LoadConfig(fs)
			if err != nil {
				return err
			}
			journal := diary.New(config, diary.WithSpecifiedFs(fs))

			since, err := time.Parse(time.DateOnly, since)
			if err != nil {
				return err
			}
			until, err := time.Parse(time.DateOnly, until)
			if err != nil {
				return err
			}

			day := since
			for !day.After(until) {
				if _, err := journal.Create(day, rewrite); err != nil {
					return err
				}
				day = day.Add(xtime.Day)
			}
			return nil
		},
	}
	flags := makeCmd.Flags()
	flags.StringVarP(&since, "since", "", since, "start date")
	flags.StringVarP(&until, "until", "", until, "end date")
	flags.BoolVarP(&rewrite, "rewrite", "", rewrite, "rewrite existing files")

	cmd.AddCommand(makeCmd)

	return cmd
}
