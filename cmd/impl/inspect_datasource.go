package impl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type datasource struct {
	name   string
	id     int
	groups []group
}

type group struct {
	name string
	id   int
	role groupRole
}

type groupRole string

var groupRoleReadonly = groupRole("Readonly")
var groupRoleFullAccess = groupRole("FullAccess")

type groupIDToNameMap map[int]string

func buildGroupMap(client redashClient) (groupIDToNameMap, error) {
	resp, err := requestGetGroups(client)
	if err != nil {
		return nil, errors.Wrap(err, "requestGetGroups")
	}

	res := groupIDToNameMap{}
	for _, g := range resp {
		res[g.ID] = g.Name
	}

	return res, nil
}

func inspectDataSource(client redashClient, dataSourceID int) error {
	groupMap, err := buildGroupMap(client)
	if err != nil {
		return errors.Wrap(err, "buildGroupMap")
	}

	ds, err := buildDatasource(client, groupMap, dataSourceID)
	if err != nil {
		return errors.Wrap(err, "buildDatasource")
	}

	explainDatasource(ds, 0)
	return nil
}

func buildDatasource(client redashClient, groupMap groupIDToNameMap, id int) (ds datasource, err error) {
	res, err := requestGetDataSource(client, id)
	if err != nil {
		return ds, errors.Wrap(err, "requestGetDataSource")
	}

	ds.name = res.Name
	ds.id = res.ID

	for groupID, readonly := range res.Groups {
		gID, err := strconv.Atoi(groupID)
		if err != nil {
			return ds, errors.Wrap(err, "strconv.Atoi")
		}

		g := group{
			name: groupMap[gID],
			id:   gID,
		}

		if readonly {
			g.role = groupRoleReadonly
		} else {
			g.role = groupRoleFullAccess
		}

		ds.groups = append(ds.groups, g)
	}

	return ds, nil
}

func explainDatasource(ds datasource, indent int) {
	fmt.Printf("%sID %d datasource is: %s\n", strings.Repeat("\t", indent), ds.id, ds.name)

	for _, g := range ds.groups {
		explainGroup(g, indent+1)
	}
}

func explainGroup(g group, indent int) {
	fmt.Printf("%s* %s (%s)\n", strings.Repeat("\t", indent), g.name, g.role)
}
