package impl

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

func InspectQueryCmd(client inspectRedashClient, queryID int, acceptReadonly bool) error {
	if err := explainQuery(client, queryID, 0); err != nil {
		return xerrors.Errorf("explainQuery: %+w", err)
	}

	dataSourceID, err := getDataSourceIDByQuery(client, queryID)
	if err != nil {
		return xerrors.Errorf("getDataSourceIDByQuery: %+w", err)
	}

	if err := explainDataSource(client, dataSourceID, 1); err != nil {
		return xerrors.Errorf("explainDataSource: %+w", err)
	}

	ids, err := getGroupIDsByDataSource(client, dataSourceID, acceptReadonly)
	if err != nil {
		return xerrors.Errorf("getGroupIDsByDataSource: %+w", err)
	}

	if err := explainGroups(client, ids, 2); err != nil {
		return xerrors.Errorf("explainGroups: %+w", err)
	}

	return nil
}

func InspectDataSourceCmd(client inspectRedashClient, dataSourceID int, acceptReadonly bool) error {
	if err := explainDataSource(client, dataSourceID, 0); err != nil {
		return xerrors.Errorf("explainDataSource: %+w", err)
	}

	ids, err := getGroupIDsByDataSource(client, dataSourceID, acceptReadonly)
	if err != nil {
		return xerrors.Errorf("getGroupIDsByDataSource: %+w", err)
	}

	if err := explainGroups(client, ids, 1); err != nil {
		return xerrors.Errorf("explainGroups: %+w", err)
	}

	return nil
}

func InspectDashboardCmd(client inspectRedashClient, slug string, acceptReadonly bool) error {
	if err := explainDashboard(client, slug, 0); err != nil {
		return xerrors.Errorf("explainDashboard: %+w", err)
	}

	queryIDs, err := getQueryIDsByDashboard(client, slug)
	if err != nil {
		return xerrors.Errorf("getQueryIDsByDashboard: %+w", err)
	}

	for _, qID := range queryIDs {
		if err := explainQuery(client, qID, 1); err != nil {
			return xerrors.Errorf("explainQuery: %+w", err)
		}

		dataSourceID, err := getDataSourceIDByQuery(client, qID)
		if err != nil {
			return xerrors.Errorf("getDataSourceIDByQuery: %+w", err)
		}

		if err := explainDataSource(client, dataSourceID, 2); err != nil {
			return xerrors.Errorf("explainDataSource: %+w", err)
		}

		ids, err := getGroupIDsByDataSource(client, dataSourceID, acceptReadonly)
		if err != nil {
			return xerrors.Errorf("getGroupIDsByDataSource: %+w", err)
		}

		if err := explainGroups(client, ids, 3); err != nil {
			return xerrors.Errorf("explainGroups: %+w", err)
		}
	}
	return nil
}

func getDataSourceIDByQuery(client inspectRedashClient, queryID int) (int, error) {
	raw, err := client.GetQuery(queryID)
	if err != nil {
		return 0, xerrors.Errorf("client.GetQuery: %+w", err)
	}

	var resp struct {
		Id           int    `json:"id"`
		Name         string `json:"name"`
		DataSourceId int    `json:"data_source_id"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return 0, xerrors.Errorf("json.Unmarshal: %+w", err)
	}

	return resp.DataSourceId, nil
}

func getGroupIDsByDataSource(client inspectRedashClient, dataSourceID int, acceptReadOnly bool) ([]int, error) {
	raw, err := client.GetDataSource(dataSourceID)
	if err != nil {
		return nil, xerrors.Errorf("client.GetDataSource: %+w", err)
	}
	var resp struct {
		Groups map[string]bool `json:"groups"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, xerrors.Errorf("json.Unmarshal: %+w", err)
	}

	var res []int

	for k, readonly := range resp.Groups {
		groupID, err := strconv.Atoi(k)
		if err != nil {
			return nil, xerrors.Errorf("strconv.Atoi: %+w", err)
		}

		if readonly {
			if acceptReadOnly {
				res = append(res, groupID)
			}
		} else {
			res = append(res, groupID)
		}
	}

	return res, nil
}

func getGroupIDToName(client inspectRedashClient) (map[int]string, error) {
	raw, err := client.GetGroups()
	if err != nil {
		return nil, xerrors.Errorf("client.GetGroups: %+w", err)
	}
	var resp []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, xerrors.Errorf("json.Unmarshal: %+w", err)
	}

	res := make(map[int]string)
	for _, g := range resp {
		res[g.Id] = g.Name
	}

	return res, nil
}

func getQueryIDsByDashboard(client inspectRedashClient, dashboardSlug string) ([]int, error) {
	raw, err := client.GetDashboard(dashboardSlug)
	if err != nil {
		return nil, xerrors.Errorf("client.GetDashboard: %+w", err)
	}

	var resp struct {
		Widgets []struct {
			Visualization *struct {
				Query struct {
					Id int `json:"id"`
				} `json:"query"`
			} `json:"visualization"`
		} `json:"widgets"`
	}

	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, xerrors.Errorf("json.Unmarshal: %+w", err)
	}

	var ids []int
	for _, w := range resp.Widgets {
		if w.Visualization == nil {
			continue
		}

		ids = append(ids, w.Visualization.Query.Id)
	}

	return ids, nil
}

func explainDashboard(client inspectRedashClient, dashboardSlug string, indent int) error {
	raw, err := client.GetDashboard(dashboardSlug)
	if err != nil {
		return xerrors.Errorf("client.GetDashboard: %+w", err)
	}

	var resp struct {
		Slug string `json:"slug"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return xerrors.Errorf("json.Unmarshal: %+w", err)
	}

	fmt.Printf("%s slug `%s` dashboard is: %s\n", strings.Repeat("\t", indent), dashboardSlug, resp.Name)
	return nil
}

func explainQuery(client inspectRedashClient, queryID int, indent int) error {
	raw, err := client.GetQuery(queryID)
	if err != nil {
		return xerrors.Errorf("client.GetQuery: %+w", err)
	}

	var resp struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return xerrors.Errorf("json.Unmarshal: %+w", err)
	}

	fmt.Printf("%s ID %d query is: %s\n", strings.Repeat("\t", indent), queryID, resp.Name)

	return nil
}

func explainDataSource(client inspectRedashClient, dataSourceID int, indent int) error {
	raw, err := client.GetDataSource(dataSourceID)
	if err != nil {
		return xerrors.Errorf("client.GetDataSource: %+w", err)
	}
	var resp struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return xerrors.Errorf("json.Unmarshal: %+w", err)
	}

	fmt.Printf("%s ID %d datasource is: %s\n", strings.Repeat("\t", indent), dataSourceID, resp.Name)

	return nil
}

func explainGroups(client inspectRedashClient, ids []int, indent int) error {
	groupMap, err := getGroupIDToName(client)
	if err != nil {
		return xerrors.Errorf("getGroupIDToName: %+w", err)
	}

	for _, id := range ids {
		fmt.Printf("%s avaliable group: `%s`\n", strings.Repeat("\t", indent), groupMap[id])
	}

	return nil
}

type inspectRedashClient interface {
	GetQuery(id int) ([]byte, error)
	GetDataSource(id int) ([]byte, error)
	GetGroups() ([]byte, error)
	GetDashboard(id string) ([]byte, error)
}
