package sparkle

import (
	"bufio"

	"github.com/spf13/cobra"

	"go.octolab.org/ecosystem/sparkle/internal/service/tact"
)

func Tact() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tact",
	}

	cmd.AddCommand(
		Logbook(),
	)

	return cmd
}

func Logbook() *cobra.Command {
	cmd := &cobra.Command{
		Use: "logbook",
	}

	calculate := &cobra.Command{
		Use:  "calculate",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			journal := new(tact.Logbook)
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
			cmd.Println(journal)
			return nil
		},
	}

	cmd.AddCommand(calculate)

	return cmd
}
