package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

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

	return mainAddCmd(client, users, groups)
}

func mainAddCmd(client redashClient, users, groups []string) error {

	fmt.Printf("users:   %s\n", users)
	fmt.Printf("groups:  %s\n", groups)
	fmt.Printf("Are you sure? [y/n]")

	reader := bufio.NewReader(os.Stdin)
	res, err := reader.ReadString('\n')

	if err != nil || strings.TrimSpace(res) != "y" {
		fmt.Println("Abort.")
		return nil
	}

	groupIds := make([]int, len(groups))
	for i, g := range groups {
		id, err := findGroupID(client, g)
		if err != nil {
			return xerrors.Errorf("%+w", err)
		}
		groupIds[i] = id
	}

	userIds := make([]int, len(users))
	for i, u := range users {
		id, err := findUserID(client, u)
		if err != nil {
			return xerrors.Errorf("%+w", err)
		}
		userIds[i] = id
	}

	for i, g := range groupIds {
		for j, u := range userIds {
			_, err := client.AddMember(g, u)
			if err != nil {
				return xerrors.Errorf("client.Addmember: %+w", err)
			}

			fmt.Printf("Added %s to %s\n", users[j], groups[i])
		}
	}

	return nil
}

type redashClient interface {
	SearchUser(q string) ([]byte, error)
	GetGroups() ([]byte, error)
	AddMember(groupId, userId int) ([]byte, error)
}

func findGroupID(client redashClient, groupName string) (int, error) {
	raw, err := client.GetGroups()
	if err != nil {
		return -1, err
	}

	var resp []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return -1, err
	}

	for _, g := range resp {
		if g.Name == groupName {
			return g.Id, nil
		}
	}

	return -1, xerrors.Errorf("group: %s not found", groupName)
}

func findUserID(client redashClient, userEmail string) (int, error) {
	raw, err := client.SearchUser(userEmail)
	if err != nil {
		return -1, err
	}

	var resp struct {
		Results []struct {
			Id    int    `json:"id"`
			Email string `json:"email"`
		} `json:"results"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return -1, err
	}

	if len(resp.Results) == 0 {
		return -1, xerrors.Errorf("user: %s not found", userEmail)
	}

	user := resp.Results[0]
	if user.Email != userEmail {
		return -1, xerrors.Errorf("user: %s not found", userEmail)
	}

	return user.Id, nil
}
