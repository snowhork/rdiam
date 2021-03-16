package redash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/xerrors"
)

type Client struct {
	apiKey   string
	endpoint string
}

func NewClient(endpoint, apiKey string) *Client {
	return &Client{apiKey: apiKey, endpoint: endpoint}
}

func (c *Client) get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, xerrors.Errorf("http.Get error: %+w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("ioutil.ReadAll: %+w", err)
	}

	if resp.StatusCode != 200 {
		return nil, xerrors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return body, nil
}

func (c *Client) post(url string, req io.Reader) ([]byte, error) {
	resp, err := http.Post(url, "application/json", req)
	if err != nil {
		return nil, xerrors.Errorf("http.Post error: %+w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("ioutil.ReadAll: %+w", err)
	}

	if resp.StatusCode != 200 {
		return nil, xerrors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return body, nil
}

func (c *Client) SearchUser(q string) ([]byte, error) {
	val := url.Values{}
	val.Add("q", q)
	val.Add("api_key", c.apiKey)

	return c.get(fmt.Sprintf("%s/api/users?%s", c.endpoint, val.Encode()))
}

func (c *Client) GetGroups() ([]byte, error) {
	val := url.Values{}
	val.Add("api_key", c.apiKey)

	return c.get(fmt.Sprintf("%s/api/groups?%s", c.endpoint, val.Encode()))
}

func (c *Client) GetDataSource(id int) ([]byte, error) {
	val := url.Values{}
	val.Add("api_key", c.apiKey)

	return c.get(fmt.Sprintf("%s/api/data_sources/%d?%s", c.endpoint, id, val.Encode()))
}

func (c *Client) GetQuery(id int) ([]byte, error) {
	val := url.Values{}
	val.Add("api_key", c.apiKey)

	return c.get(fmt.Sprintf("%s/api/queries/%d?%s", c.endpoint, id, val.Encode()))
}

func (c *Client) GetDashboard(id string) ([]byte, error) {
	val := url.Values{}
	val.Add("api_key", c.apiKey)

	return c.get(fmt.Sprintf("%s/api/dashboards/%s?%s", c.endpoint, id, val.Encode()))
}

func (c *Client) AddMember(groupID, userID int) ([]byte, error) {
	val := url.Values{}
	val.Add("api_key", c.apiKey)

	body := struct {
		UserID int `json:"user_id"`
	}{UserID: userID}

	raw, err := json.Marshal(body)
	if err != nil {
		return nil, xerrors.Errorf("json.Marshal: %+w", err)
	}

	return c.post(fmt.Sprintf("%s/api/groups/%d/members?%s", c.endpoint, groupID, val.Encode()), bytes.NewBuffer(raw))

}
