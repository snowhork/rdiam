package impl

import (
	"fmt"

	"golang.org/x/xerrors"
)

func AddCmd(client redashClient, users, groups []string) error {
	groupIds := make([]int, len(groups))
	for i, g := range groups {
		id, err := findGroupID(client, g)
		if err != nil {
			return xerrors.Errorf("findGroupID: %+w", err)
		}
		groupIds[i] = id
	}

	userIds := make([]int, len(users))
	for i, u := range users {
		id, err := findUserID(client, u)
		if err != nil {
			return xerrors.Errorf("findUserID: %+w", err)
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
