package impl

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type acl struct {
	userID   int
	userName string
}

type query struct {
	queryName  string
	queryID    int
	queryACL   []acl
	datasource datasource
}

func inspectQuery(client redashClient, queryID int, o Options) error {
	groupMap, err := buildGroupMap(client)
	if err != nil {
		return errors.Wrap(err, "buildGroupMap")
	}

	q, err := buildQuery(client, groupMap, queryID)
	if err != nil {
		return errors.Wrap(err, "buildQuery")
	}

	explainQuery(q, 0, o)
	return nil
}

func buildQuery(client redashClient, groupMap groupIDToNameMap, queryID int) (q query, err error) {
	res, err := requestGetQuery(client, queryID)
	if err != nil {
		return q, errors.Wrap(err, "requestGetQuery")
	}

	q.queryName = res.Name
	q.queryID = res.ID

	resACL, err := requestGetQueryACL(client, queryID)
	if err != nil {
		return q, errors.Wrap(err, "requestGetQueryACL")
	}
	for _, m := range resACL.Modify {
		q.queryACL = append(q.queryACL, acl{m.ID, m.Name})
	}

	ds, err := buildDatasource(client, groupMap, res.DataSourceID)
	if err != nil {
		return q, errors.Wrap(err, "buildDatasource")
	}

	q.datasource = ds
	return q, nil
}

func explainQuery(q query, indent int, o Options) {
	fmt.Printf("%sID %d query is: %s\n", strings.Repeat("\t", indent), q.queryID, q.queryName)
	explainDatasource(q.datasource, indent+1)
	if o.ShowQueryModifyPermission && len(q.queryACL) > 0 {
		fmt.Printf("%sUsers with modify permission:\n", strings.Repeat("\t", indent+1))
		for _, a := range q.queryACL {
			fmt.Printf("%s* UserID: %d(%s)\n", strings.Repeat("\t", indent+2), a.userID, a.userName)
		}
	}
}
