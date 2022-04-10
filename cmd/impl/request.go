package impl

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type redashClient interface {
	SearchUser(q string) ([]byte, error)
	GetGroups() ([]byte, error)
	GetQuery(id int) ([]byte, error)
	GetQueryACL(id int) ([]byte, error)
	GetDataSource(id int) ([]byte, error)
	GetDashboard(id string) ([]byte, error)
	AddMember(groupID, userID int) ([]byte, error)
}

type responseSearchUser struct {
	Results []struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	} `json:"results"`
}

func requestSearchUser(client redashClient, q string) (resp responseSearchUser, err error) {
	raw, err := client.SearchUser(q)
	if err != nil {
		return resp, errors.Wrap(err, "client.SearchUser")
	}

	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, errors.Wrap(err, "json.Unmarshal")
	}

	return resp, nil
}

type responseGetGroups []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func requestGetGroups(client redashClient) (resp responseGetGroups, err error) {
	raw, err := client.GetGroups()
	if err != nil {
		return resp, errors.Wrap(err, "client.SearchUser")
	}

	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, errors.Wrap(err, "json.Unmarshal")
	}

	return resp, nil
}

type responseGetQuery struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DataSourceID int    `json:"data_source_id"`
}

func requestGetQuery(client redashClient, id int) (resp responseGetQuery, err error) {
	raw, err := client.GetQuery(id)
	if err != nil {
		return resp, errors.Wrap(err, "client.GetQuery")
	}

	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, errors.Wrap(err, "json.Unmarshal")
	}

	return resp, nil
}

type responseGetQueryACL struct {
	Modify []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"modify"`
}

func requestGetQueryACL(client redashClient, id int) (resp responseGetQueryACL, err error) {
	raw, err := client.GetQueryACL(id)
	if err != nil {
		return resp, errors.Wrap(err, "client.GetQueryACL")
	}

	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, errors.Wrap(err, "json.Unmarshal")
	}

	return resp, nil
}

type responseGetDataSource struct {
	ID     int             `json:"id"`
	Name   string          `json:"name"`
	Groups map[string]bool `json:"groups"`
}

func requestGetDataSource(client redashClient, id int) (resp responseGetDataSource, err error) {
	raw, err := client.GetDataSource(id)
	if err != nil {
		return resp, errors.Wrap(err, "client.GetDataSource")
	}

	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, errors.Wrap(err, "json.Unmarshal")
	}

	return resp, nil
}

type responseGetDashboard struct {
	Slug    string `json:"slug"`
	Name    string `json:"name"`
	Widgets []struct {
		Visualization *struct {
			Query struct {
				ID int `json:"id"`
			} `json:"query"`
		} `json:"visualization"`
	} `json:"widgets"`
}

func requestGetDashboard(client redashClient, slug string) (resp responseGetDashboard, err error) {
	raw, err := client.GetDashboard(slug)
	if err != nil {
		return resp, errors.Wrap(err, "client.GetDashboard")
	}

	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, errors.Wrap(err, "json.Unmarshal")
	}

	return resp, nil
}
