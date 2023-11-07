package sparkle

import (
	"time"

	"github.com/nleeper/goment"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	xtime "go.octolab.org/time"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
	diary "go.octolab.org/ecosystem/sparkle/internal/plugins/obsidian/daily-notes"
	"go.octolab.org/ecosystem/sparkle/internal/plugins/obsidian/periodic-notes"
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

			// TODO:deps:refactoring improve configuration
			cnf, err := periodic.LoadConfig(fs)
			if err != nil {
				return err
			}
			transformers := []func(*goment.Goment) func(*markdown.Document){
				periodic.LinkRelatives(cnf.Weekly),
				periodic.LinkRelatives(cnf.Monthly),
				periodic.LinkRelatives(cnf.Quarterly),
				periodic.LinkRelatives(cnf.Yearly),
			}

			day := from
			for !day.After(to) {
				if _, err := journal.Create(day, transformers...); err != nil {
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
