package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"

	"github.com/snowhork/rdiam/cmd/impl"
)

func init() {
	inspectCmd := newInspectCmd()
	rootCmd.AddCommand(inspectCmd)

	inspectCmd.AddCommand(newInspectQueryCmd())
	inspectCmd.AddCommand(newInspectDataSourceCmd())
	inspectCmd.AddCommand(newInspectDashboardCmd())
}

func newInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect [query|datasource|dashboard]",
		Short: "display required group for query or datasource or dashboard",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return xerrors.New("query id is required")
			}

			return nil
		},
	}

	return cmd
}

func newInspectQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query 123",
		Short: "display required group for datasource of the query",
		Long: `
For example:

rdiam inspect query 123`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return xerrors.New("query id is required")
			}

			return nil
		},
		RunE: runInspectQueryCmd,
	}

	return cmd
}

func runInspectQueryCmd(cmd *cobra.Command, args []string) error {
	queryID, err := strconv.Atoi(args[0])
	if err != nil {
		return xerrors.Errorf("queryID must be integer: %+w", err)
	}

	return impl.InspectQueryCmd(globalClient, queryID, false)
}

func newInspectDataSourceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datasource 12345",
		Short: "display required group for the datasource",
		Long: `
For example:

rdiam inspect datasource 12345`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return xerrors.New("query id is required")
			}

			return nil
		},
		RunE: runInspectDataSourceCmd,
	}

	return cmd
}

func runInspectDataSourceCmd(cmd *cobra.Command, args []string) error {
	dataSourceID, err := strconv.Atoi(args[0])
	if err != nil {
		return xerrors.Errorf("dataSourceID must be integer: %+w", err)
	}

	return impl.InspectDataSourceCmd(globalClient, dataSourceID, false)
}

func newInspectDashboardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboard board-name",
		Short: "display queries and required group for each query",
		Long: `
For example:

rdiam inspect dashboard board-name`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return xerrors.New("query id is required")
			}

			return nil
		},
		RunE: runInspectDashboardCmd,
	}

	return cmd
}

func runInspectDashboardCmd(cmd *cobra.Command, args []string) error {
	dashboardID := args[0]
	return impl.InspectDashboardCmd(globalClient, dashboardID, false)
}
