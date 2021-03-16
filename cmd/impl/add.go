package impl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/xerrors"
)

func AddCmd(client redashClient, users, groups []string) error {
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

func findGroupID(client redashClient, groupName string) (int, error) {
	resp, err := requestGetGroups(client)
	if err != nil {
		return -1, xerrors.Errorf("requestGetGroups: %+w", err)
	}

	for _, g := range resp {
		if g.Name == groupName {
			return g.ID, nil
		}
	}

	return -1, xerrors.Errorf("group: %s not found", groupName)
}

func findUserID(client redashClient, userEmail string) (int, error) {
	resp, err := requestSearchUser(client, userEmail)
	if err != nil {
		return -1, xerrors.Errorf("requestSearchUser: %+w", err)
	}

	if len(resp.Results) == 0 {
		return -1, xerrors.Errorf("user: %s not found", userEmail)
	}

	user := resp.Results[0]
	if user.Email != userEmail {
		return -1, xerrors.Errorf("user: %s not found", userEmail)
	}

	return user.ID, nil
}
