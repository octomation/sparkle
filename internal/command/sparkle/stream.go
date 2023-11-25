package sparkle

import (
	"bufio"
	"os"
	"time"

	"github.com/nleeper/goment"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
	xtime "go.octolab.org/ecosystem/sparkle/internal/pkg/x/time"
	diary "go.octolab.org/ecosystem/sparkle/internal/plugins/obsidian/daily-notes"
	"go.octolab.org/ecosystem/sparkle/internal/plugins/obsidian/periodic-notes"
	"go.octolab.org/ecosystem/sparkle/internal/service/tact"
)

func Stream() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "stream category",
		Args: cobra.NoArgs,
	}
	cmd.AddCommand(
		Diary(),
		Logbook(),
		Plans(),
	)
	return cmd
}

func Diary() *cobra.Command {
	fs := afero.NewOsFs()
	cmd := &cobra.Command{
		Use:  "diary command",
		Args: cobra.NoArgs,
	}

	var (
		since   = time.Now().Format(time.DateOnly)
		until   = time.Now().Add(xtime.Week).Format(time.DateOnly)
		next    = false
		rewrite = false
	)
	makeCmd := &cobra.Command{
		Use:  "make [--since=YYYY-MM-DD] [--until=YYYY-MM-DD] [--rewrite]",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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
			transformers := []func(*goment.Goment) markdown.Transformer{
				periodic.LinkRelatives(cnf.Weekly),
				periodic.LinkRelatives(cnf.Monthly),
				periodic.LinkRelatives(cnf.Quarterly),
				periodic.LinkRelatives(cnf.Yearly),
				periodic.LinkSiblings(cnf.Daily, periodic.LookupDays),
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
	flags.BoolVarP(&rewrite, "rewrite", "", rewrite, "rewrite existing files")

	cmd.AddCommand(makeCmd)

	return cmd
}

func Logbook() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "logbook command",
		Args: cobra.NoArgs,
	}

	calculate := &cobra.Command{
		Use:  "calculate {stdin}",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			journal := tact.NewLinearJournal()
			scanner := bufio.NewScanner(cmd.InOrStdin())
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				if err := journal.Log(scanner.Text()); err != nil {
					return err
				}
			}
			if err := scanner.Err(); err != nil {
				return err
			}

			// temporary hack
			cmd.SetOut(os.Stdout)
			cmd.Println(journal)
			return nil
		},
	}

	cmd.AddCommand(calculate)

	return cmd
}

func Plans() *cobra.Command {
	fs := afero.NewOsFs()
	cmd := &cobra.Command{
		Use:  "plans command",
		Args: cobra.NoArgs,
	}

	var (
		week    bool
		month   bool
		quarter bool
		year    bool
		next    bool
	)
	makeCmd := &cobra.Command{
		Use:  "make [--next] {--week --month --quarter --year}",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := periodic.LoadConfig(fs)
			if err != nil {
				return err
			}
			planner := periodic.New(
				config,
				periodic.WithSpecifiedFs(fs),
				periodic.WithTransformers(
					periodic.UpdateAliases(),
					// we don't know cascade relations
					// and have to link all relatives
					periodic.LinkRelatives(config.Weekly),
					periodic.LinkRelatives(config.Monthly),
					periodic.LinkRelatives(config.Quarterly),
					periodic.LinkRelatives(config.Yearly),
				),
			)

			shift := func(ref time.Time, fn func(time.Time) time.Time) time.Time {
				if next {
					return fn(ref)
				}
				return ref
			}

			now := time.Now()
			if week {
				_, err := planner.Week(
					shift(now, xtime.NextWeek),
					periodic.LinkSiblings(config.Weekly, periodic.LookupWeeks),
				)
				if err != nil {
					return err
				}
			}
			if month {
				_, err := planner.Month(
					shift(now, xtime.NextMonth),
					periodic.LinkSiblings(config.Monthly, periodic.LookupMonths),
				)
				if err != nil {
					return err
				}
			}
			if quarter {
				_, err := planner.Quarter(
					shift(now, xtime.NextQuarter),
					periodic.LinkSiblings(config.Quarterly, periodic.LookupQuarters),
				)
				if err != nil {
					return err
				}
			}
			if year {
				_, err := planner.Year(
					shift(now, xtime.NextYear),
					periodic.LinkSiblings(config.Yearly, periodic.LookupYears),
				)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	flags := makeCmd.Flags()
	flags.BoolVar(&week, "week", false, "make plans for the week")
	flags.BoolVar(&month, "month", false, "make plans for the month")
	flags.BoolVar(&quarter, "quarter", false, "make plans for the quarter")
	flags.BoolVar(&year, "year", false, "make plans for the year")
	flags.BoolVar(&next, "next", false, "make plans for the next period, e.g., next week, etc.")

	cmd.AddCommand(makeCmd)

	return cmd
}
