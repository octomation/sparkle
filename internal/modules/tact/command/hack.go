package command

import (
	"fmt"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	xtime "go.octolab.org/time"

	diary "go.octolab.org/ecosystem/sparkle/internal/plugins/obsidian/daily-notes"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tact",
		RunE: func(cmd *cobra.Command, args []string) error {
			fs := afero.NewOsFs()
			config, err := diary.LoadConfig(fs)
			if err != nil {
				return err
			}
			journal := diary.New(config, diary.WithSpecifiedFs(fs))

			since, err := time.Parse("2006-01-02", "2023-10-30")
			if err != nil {
				return err
			}
			until, err := time.Parse("2006-01-02", "2023-11-06")
			if err != nil {
				return err
			}

			day := since
			for day.Before(until) {
				entry, err := journal.Create(day, true)
				if err != nil {
					return err
				}
				fmt.Println(entry.Path())
				day = day.Add(xtime.Day)
			}

			return nil
		},
	}

	return cmd
}
