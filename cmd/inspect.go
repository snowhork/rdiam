package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"

	"github.com/snowhork/rdiam/cmd/impl"
)

var (
	o impl.Options
)

func init() {
	rootCmd.AddCommand(newInspectCmd())
}

func newInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect https://redash.yourdomain.com/queries/12345",
		Short: "display required group for query or datasource or dashboard",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return xerrors.New("only 1 url must be specified")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var arg string
			if len(args) == 0 {
				fmt.Printf("Type target URL (e.g. https://redash.yourdomain.com/queries/12345): ")

				reader := bufio.NewReader(os.Stdin)
				res, err := reader.ReadString('\n')

				if err != nil {
					return err
				}

				arg = strings.TrimSpace(res)
			} else {
				arg = args[0]
			}
			impl.SetInspectOptions(o)

			return impl.InspectCmd(globalClient, arg)
		},
	}
	cmd.Flags().BoolVarP(&o.ShowQueryModifyPermission, "show-permission", "p", false, "Show query modify permission")

	return cmd
}
