package impl

type mockClient struct {
	respSearchUser,
	respGetGroups,
	respGetQuery,
	respDataSource,
	respDashboard,
	respAddMember []byte
}

func (c *mockClient) SearchUser(q string) ([]byte, error) {
	return c.respSearchUser, nil
}

func (c *mockClient) GetGroups() ([]byte, error) {
	return c.respGetGroups, nil
}

func (c *mockClient) GetQuery(id int) ([]byte, error) {
	return c.respGetQuery, nil
}

func (c *mockClient) GetDataSource(id int) ([]byte, error) {
	return c.respDataSource, nil
}

func (c *mockClient) GetDashboard(id string) ([]byte, error) {
	return c.respDashboard, nil
}

func (c *mockClient) AddMember(groupID, userID int) ([]byte, error) {
	return c.respAddMember, nil
}
