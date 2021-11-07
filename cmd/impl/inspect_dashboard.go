package impl

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

type dashboard struct {
	name    string
	slug    string
	queries []query
}

func inspectDashboard(client redashClient, slug string, o Options) error {
	groupMap, err := buildGroupMap(client)
	if err != nil {
		return xerrors.Errorf("buildGroupMap: %+w", err)
	}

	db, err := buildDashboard(client, groupMap, slug)
	if err != nil {
		return xerrors.Errorf("buildQuery: %+w", err)
	}

	explainDashboard(db, 0, o)
	return nil
}

func buildDashboard(client redashClient, groupMap groupIDToNameMap, slug string) (db dashboard, err error) {
	res, err := requestGetDashboard(client, slug)
	if err != nil {
		return db, xerrors.Errorf("requestGetDashboard: %+w", err)
	}

	db.slug = res.Slug
	db.name = res.Name

	for _, w := range res.Widgets {
		if w.Visualization == nil {
			continue
		}

		q, err := buildQuery(client, groupMap, w.Visualization.Query.ID)
		if err != nil {
			return db, xerrors.Errorf("buildQuery: %+w", err)
		}
		db.queries = append(db.queries, q)
	}

	return db, nil
}

func explainDashboard(db dashboard, indent int, o Options) {
	fmt.Printf("%sID %s dashboard is: %s\n", strings.Repeat("\t", indent), db.slug, db.name)
	for _, q := range db.queries {
		explainQuery(q, indent+1, o)
	}
}
