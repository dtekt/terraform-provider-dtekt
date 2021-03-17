package dtekt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetTests(account_id string) (*RunConfiguration, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/runconfiguration/", c.HostURL, account_id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	rc := RunConfiguration{}
	err = json.Unmarshal(body, &rc)

	if err != nil {
		return nil, err
	}

	return &rc, nil
}

func (c *Client) GetTest(testUuid string, account_id string) (*Run, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/runconfiguration/%s", c.HostURL, account_id, testUuid), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := Run{}
	err = json.Unmarshal(body, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) CreateTest(test Run, account_id string) (*Run, error) {
	test.Count = 0
	rb, err := json.Marshal(test)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/account/%s/runconfiguration/", c.HostURL, account_id), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	test_created := Run{}
	err = json.Unmarshal(body, &test_created)
	if err != nil {
		return nil, err
	}

	return &test_created, nil
}

func (c *Client) DeleteTest(testUuid string, account_id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/runconfiguration/%s", c.HostURL, account_id, testUuid), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
