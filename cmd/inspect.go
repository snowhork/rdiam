package cmd

import (
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"

	"github.com/snowhork/rdiam/cmd/impl"
)

func init() {
	rootCmd.AddCommand(newInspectCmd())
}

func newInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect https://redash.yourdomain.com/queries/12345",
		Short: "display required group for query or datasource or dashboard",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return xerrors.New("query id is required")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return impl.InspectCmd(globalClient, args[0])
		},
	}

	return cmd
}
