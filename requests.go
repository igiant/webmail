package webmail

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

type ClientConnection struct {
	Config *Config
	Token  *string
	client *http.Client
}

func (c *Config) NewConnection() (*ClientConnection, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	connection := ClientConnection{
		Config: c,
		client: client,
	}
	return &connection, nil
}

func newClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &http.Client{Jar: jar}, nil
}

func (c *ClientConnection) CallRaw(method string, params interface{}) ([]byte, error) {
	buffer, err := marshal(c.Config.getID(), method, c.Token, params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.Config.url, bytes.NewBuffer(buffer))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "ApiApplication/json-rpc")
	if c.Token != nil {
		req.Header.Add("X-Token", *c.Token)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = checkError(data); err != nil {
		return nil, err
	}
	return data, nil
}

func addMissedParametersToSearchQuery(query SearchQuery) SearchQuery {
	if query.Fields == nil {
		query.Fields = []string{}
	}
	if query.Conditions == nil {
		query.Conditions = SubConditionList{}
	}
	if query.Combining == "" {
		query.Combining = Or
	}
	return query
}
