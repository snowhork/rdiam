package impl

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Options struct {
	ShowQueryModifyPermission bool
}

func InspectCmd(client redashClient, raw string, o Options) error {
	u, err := url.Parse(raw)
	if err != nil {
		return errors.Wrap(err, "url.Parse")
	}

	path := strings.Split(u.Path, "/")
	if len(path) < 3 {
		return errors.Errorf("unable to inspect the url: %q", raw)
	}
	resource, id := path[1], path[2]

	switch resource {
	case "queries":
		qID, err := strconv.Atoi(id)
		if err != nil {
			return errors.Wrap(err, "query id must be integer")
		}
		return inspectQuery(client, qID, o)

	case "data_sources":
		dID, err := strconv.Atoi(id)
		if err != nil {
			return errors.Wrap(err, "data_source id must be integer")
		}
		return inspectDataSource(client, dID)

	case "dashboard":
		return inspectDashboard(client, id, o)
	}
	return errors.Errorf("unknown resource type: %q", resource)
}
