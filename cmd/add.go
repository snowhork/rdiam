package cmd

import (
	"os"

	"github.com/snowhork/rdiam/cmd/impl"

	"github.com/snowhork/rdiam/pkg/redash"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

func init() {
	rootCmd.AddCommand(newAddCmd())
}

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add -u user1@email.com,user2@email.com -g group1,group2",
		Short: "add some users to some groups",
		Long: `add some users to some groups
For example:

rdiam add -u user1@email.com,user2@email.com -g group1,group2`,
		RunE: runAddCmd,
	}

	cmd.Flags().StringSliceP("users", "u", []string{}, "Specify user email(s)")
	cmd.Flags().StringSliceP("groups", "g", []string{}, "Specify group(s)")

	if err := cmd.MarkFlagRequired("users"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("groups"); err != nil {
		panic(err)
	}

	return cmd
}

func runAddCmd(cmd *cobra.Command, args []string) error {
	users, err := cmd.Flags().GetStringSlice("users")
	if err != nil {
		return xerrors.Errorf("failed to parse users flag: %+w", err)
	}

	groups, err := cmd.Flags().GetStringSlice("groups")
	if err != nil {
		return xerrors.Errorf("failed to parse groups flag: %+w", err)
	}
	client := redash.NewClient(os.Getenv("REDASH_ENDPOINT"), os.Getenv("REDASH_API_KEY"))

	return impl.AddCmd(client, users, groups)
}
