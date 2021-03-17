package dtekt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// func (c *Client) GetAccounts() (*RunConfiguration, error) {
// 	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/ZcVSEv9a4B9FnFC8Hb7dM6/runconfiguration", c.HostURL), nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	body, err := c.doRequest(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	rc := RunConfiguration{}
// 	err = json.Unmarshal(body, &rc)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &rc, nil
// }

func (c *Client) GetAccount(accountId string) (*Account, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/", c.HostURL, accountId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := Account{}
	err = json.Unmarshal(body, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) CreateAccount(account Account) (*Account, error) {
	rb, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	account_created := Account{}
	err = json.Unmarshal(body, &account_created)
	if err != nil {
		return nil, err
	}

	return &account_created, nil
}

func (c *Client) DeleteAccount(accountId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/", c.HostURL, accountId), nil)
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
