package impl

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"
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

func inspectQuery(client redashClient, queryID int) error {
	groupMap, err := buildGroupMap(client)
	if err != nil {
		return xerrors.Errorf("buildGroupMap: %+w", err)
	}

	q, err := buildQuery(client, groupMap, queryID)
	if err != nil {
		return xerrors.Errorf("buildQuery: %+w", err)
	}

	explainQuery(q, 0)
	return nil
}

func buildQuery(client redashClient, groupMap groupIDToNameMap, queryID int) (q query, err error) {
	res, err := requestGetQuery(client, queryID)
	if err != nil {
		return q, xerrors.Errorf("requestGetQuery: %+w", err)
	}

	q.queryName = res.Name
	q.queryID = res.ID

	resAcl, err := requestGetQueryAcl(client, queryID)
	if err != nil {
		return q, xerrors.Errorf("requestGetQueryAcl: %+w", err)
	}
	for _, m := range resAcl.Modify {
		q.queryACL = append(q.queryACL, acl{m.ID, m.Name})
	}

	ds, err := buildDatasource(client, groupMap, res.DataSourceID)
	if err != nil {
		return q, xerrors.Errorf("buildDatasource: %+w", err)
	}

	q.datasource = ds
	return q, nil
}

func explainQuery(q query, indent int) {
	fmt.Printf("%sID %d query is: %s\n", strings.Repeat("\t", indent), q.queryID, q.queryName)
	explainDatasource(q.datasource, indent+1)
}
