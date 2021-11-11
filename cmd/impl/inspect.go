package impl

import (
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

type Options struct {
	ShowQueryModifyPermission bool
}

func InspectCmd(client redashClient, raw string, o Options) error {
	u, err := url.Parse(raw)
	if err != nil {
		return xerrors.Errorf("unable to parse %s for url: %+w", raw, err)
	}

	path := strings.Split(u.Path, "/")
	if len(path) < 3 {
		return xerrors.Errorf("unable to inspect the url: %s", raw)
	}
	resource, id := path[1], path[2]

	switch resource {
	case "queries":
		qID, err := strconv.Atoi(id)
		if err != nil {
			return xerrors.Errorf("query id must be integer: %+w", err)
		}
		return inspectQuery(client, qID, o)

	case "data_sources":
		dID, err := strconv.Atoi(id)
		if err != nil {
			return xerrors.Errorf("data_sourcerid must be integer: %+w", err)
		}
		return inspectDataSource(client, dID)

	case "dashboard":
		return inspectDashboard(client, id, o)
	}
	return xerrors.Errorf("unknown resource type: %s", resource)
}
