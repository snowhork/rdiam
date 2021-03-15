package impl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/xerrors"
)

func AddCmd(client addRedashClient, users, groups []string) error {
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

type addRedashClient interface {
	SearchUser(q string) ([]byte, error)
	GetGroups() ([]byte, error)
	AddMember(groupID, userID int) ([]byte, error)
}

func findGroupID(client addRedashClient, groupName string) (int, error) {
	raw, err := client.GetGroups()
	if err != nil {
		return -1, err
	}

	var resp []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return -1, err
	}

	for _, g := range resp {
		if g.Name == groupName {
			return g.ID, nil
		}
	}

	return -1, xerrors.Errorf("group: %s not found", groupName)
}

func findUserID(client addRedashClient, userEmail string) (int, error) {
	raw, err := client.SearchUser(userEmail)
	if err != nil {
		return -1, err
	}

	var resp struct {
		Results []struct {
			ID    int    `json:"id"`
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

	return user.ID, nil
}
