package impl

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type dashboard struct {
	name    string
	slug    string
	queries []query
}

func inspectDashboard(client redashClient, slug string, o Options) error {
	groupMap, err := buildGroupMap(client)
	if err != nil {
		return errors.Wrap(err, "buildGroupMap")
	}

	db, err := buildDashboard(client, groupMap, slug)
	if err != nil {
		return errors.Wrap(err, "buildQuery")
	}

	explainDashboard(db, 0, o)
	return nil
}

func buildDashboard(client redashClient, groupMap groupIDToNameMap, slug string) (db dashboard, err error) {
	res, err := requestGetDashboard(client, slug)
	if err != nil {
		return db, errors.Wrap(err, "requestGetDashboard")
	}

	db.slug = res.Slug
	db.name = res.Name

	for _, w := range res.Widgets {
		if w.Visualization == nil {
			continue
		}

		q, err := buildQuery(client, groupMap, w.Visualization.Query.ID)
		if err != nil {
			return db, errors.Wrap(err, "buildQuery")
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
