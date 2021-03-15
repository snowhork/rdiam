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

type client struct {
	apiKey   string
	endpoint string
}

func NewClient(endpoint, apiKey string) *client {
	return &client{apiKey: apiKey, endpoint: endpoint}
}

func (c *client) get(url string) ([]byte, error) {
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

func (c *client) post(url string, req io.Reader) ([]byte, error) {
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

func (c *client) SearchUser(q string) ([]byte, error) {
	val := url.Values{}
	val.Add("q", q)
	val.Add("api_key", c.apiKey)

	return c.get(fmt.Sprintf("%s/api/users?%s", c.endpoint, val.Encode()))
}

func (c *client) GetGroups() ([]byte, error) {
	val := url.Values{}
	val.Add("api_key", c.apiKey)

	return c.get(fmt.Sprintf("%s/api/groups?%s", c.endpoint, val.Encode()))
}

func (c *client) AddMember(groupId, userId int) ([]byte, error) {
	val := url.Values{}
	val.Add("api_key", c.apiKey)

	body := struct {
		UserId int `json:"user_id"`
	}{UserId: userId}

	raw, err := json.Marshal(body)
	if err != nil {
		return nil, xerrors.Errorf("json.Marshal: %+w", err)
	}

	return c.post(fmt.Sprintf("%s/api/groups/%d/members?%s", c.endpoint, groupId, val.Encode()), bytes.NewBuffer(raw))

}
