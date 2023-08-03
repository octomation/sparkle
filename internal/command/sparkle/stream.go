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
		next    = false
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

			var from, to time.Time
			if next {
				last := journal.Last().Time()
				from = last.Add(xtime.Day)
				to = from.Add(xtime.Week)
			} else {
				from, err = time.Parse(time.DateOnly, since)
				if err != nil {
					return err
				}
				to, err = time.Parse(time.DateOnly, until)
				if err != nil {
					return err
				}
			}

			day := from
			for !day.After(to) {
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
	flags.BoolVarP(&next, "next", "", next, "create next seven days")
	flags.BoolVarP(&rewrite, "rewrite", "", rewrite, "rewrite existing files")

	cmd.AddCommand(makeCmd)

	return cmd
}
