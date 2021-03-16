package impl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_addCmd(t *testing.T) {
	cases := []struct {
		name   string
		users  []string
		groups []string

		valid bool
	}{
		{name: "it successes", users: []string{"user1@email.com"}, groups: []string{"my-group-3", "my-group-4"}, valid: true},
		{name: "it fails when the user is not found", users: []string{"user2@email.com"}, groups: []string{"my-group-3"}, valid: false},
		{name: "it fails when the group is not found", users: []string{"user1@email.com"}, groups: []string{"my-group-5"}, valid: false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := addCmd(&mockClient{respSearchUser: []byte(respSearchUserJSON), respGetGroups: []byte(respGetGroupJSON)}, c.users, c.groups)
			if c.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

var respSearchUserJSON = `
{
  "count": 1,
  "page": 1,
  "page_size": 20,
  "results": [
    {
      "active_at": "2021-03-16T11:12:58Z",
      "auth_type": "external",
      "created_at": "2018-08-07T01:13:40.265Z",
      "disabled_at": null,
      "email": "user1@email.com",
      "groups": [
        {
          "id": 33,
          "name": "my-group-33"
        }
      ],
      "id": 1,
      "is_disabled": false,
      "is_email_verified": true,
      "is_invitation_pending": false,
      "name": "snowhork",
      "profile_image_url": "example.com",
      "updated_at": "2021-03-16T11:13:03.856Z"
    }
  ]
}
`

var respGetGroupJSON = `
[
  {
    "created_at": "2016-05-22T09:30:53.770Z",
    "id": 1,
    "name": "admin",
    "permissions": [
      "admin",
      "super_admin"
    ],
    "type": "builtin"
  },
  {
    "created_at": "2016-05-22T09:30:53.772Z",
    "id": 2,
    "name": "default",
    "permissions": [],
    "type": "builtin"
  },
  {
    "created_at": "2016-08-22T06:37:52.138Z",
    "id": 3,
    "name": "my-group-3",
    "permissions": [
      "create_dashboard",
      "create_query",
      "edit_dashboard",
      "edit_query",
      "view_query",
      "view_source",
      "execute_query",
      "list_users",
      "schedule_query",
      "list_dashboards",
      "list_alerts",
      "list_data_sources"
    ],
    "type": "regular"
  },
  {
    "created_at": "2016-08-22T06:37:52.138Z",
    "id": 4,
    "name": "my-group-4",
    "permissions": [
      "create_dashboard",
      "create_query",
      "edit_dashboard",
      "edit_query",
      "view_query",
      "view_source",
      "execute_query",
      "list_users",
      "schedule_query",
      "list_dashboards",
      "list_alerts",
      "list_data_sources"
    ],
    "type": "regular"
  }
]
`
