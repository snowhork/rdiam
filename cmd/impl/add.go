package impl

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

func AddCmd(client redashClient, users, groups []string) error {
	groupIds := make([]int, len(groups))
	for i, g := range groups {
		id, err := findGroupID(client, g)
		if err != nil {
			return errors.Wrap(err, "findGroupID")
		}
		groupIds[i] = id
	}

	userIds := make([]int, len(users))
	var errs []error
	for i, u := range users {
		id, err := findUserID(client, u)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "findUserID"))
			continue
		}
		userIds[i] = id
	}

	if len(errs) > 0 {
		return multierr.Combine(errs...)
	}

	for i, g := range groupIds {
		for j, u := range userIds {
			_, err := client.AddMember(g, u)
			if err != nil {
				return errors.Wrap(err, "client.AddMember")
			}

			fmt.Printf("Added %s to %s\n", users[j], groups[i])
		}
	}

	return nil
}

func findGroupID(client redashClient, groupName string) (int, error) {
	resp, err := requestGetGroups(client)
	if err != nil {
		return -1, errors.Wrap(err, "requestGetGroups")
	}

	for _, g := range resp {
		if g.Name == groupName {
			return g.ID, nil
		}
	}

	return -1, errors.Errorf("group: %q not found", groupName)
}

func findUserID(client redashClient, userEmail string) (int, error) {
	resp, err := requestSearchUser(client, userEmail)
	if err != nil {
		return -1, errors.Wrap(err, "requestSearchUser")
	}

	if len(resp.Results) == 0 {
		return -1, errors.Errorf("user: %q not found", userEmail)
	}

	user := resp.Results[0]
	if user.Email != userEmail {
		return -1, errors.Errorf("user: %q not found", userEmail)
	}

	return user.ID, nil
}
