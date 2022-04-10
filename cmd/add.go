package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/snowhork/rdiam/cmd/impl"
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
	cmd.Flags().BoolP("yes", "y", false, "Automatic yes to prompts")

	if err := cmd.MarkFlagRequired("users"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("groups"); err != nil {
		panic(err)
	}

	return cmd
}

func runAddCmd(cmd *cobra.Command, _ []string) error {
	users, err := cmd.Flags().GetStringSlice("users")
	if err != nil {
		return errors.Wrap(err, "parse users flag")
	}

	groups, err := cmd.Flags().GetStringSlice("groups")
	if err != nil {
		return errors.Wrap(err, "parse groups flag")
	}

	yes, err := cmd.Flags().GetBool("yes")
	if err != nil {
		return errors.Wrap(err, "parse yes flag")
	}

	fmt.Printf("users:   %s\n", users)
	fmt.Printf("groups:  %s\n", groups)

	if !yes {
		fmt.Printf("Are you sure? [y/n]")
		reader := bufio.NewReader(os.Stdin)
		res, err := reader.ReadString('\n')

		if err != nil || strings.TrimSpace(res) != "y" {
			fmt.Println("Abort.")
			return nil
		}
	}

	return impl.AddCmd(globalClient, users, groups)
}
