package dtekt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// func (c *Client) GetAlertHandlerSlacks(account_id string) (*RunConfiguration, error) {
// 	req, err := http.NewRequest("GET", fmt.Sprintf("%s/uptime-monitors/%s/", c.HostURL, account_id), nil)
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

func (c *Client) GetAlertHandlerSlack(handlerUuid string, accountId string) (*AlertHandler, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/alerting/handlers/slack/%s/", c.HostURL, accountId, handlerUuid), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := AlertHandler{}
	err = json.Unmarshal(body, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) CreateAlertHandlerSlack(alertHandlerSlack AlertHandler, accountId string) (*AlertHandler, error) {
	rb, err := json.Marshal(alertHandlerSlack)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/alerting/handlers/slack", c.HostURL, accountId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newResource := AlertHandler{}
	err = json.Unmarshal(body, &newResource)
	if err != nil {
		return nil, err
	}

	return &newResource, nil
}

func (c *Client) UpdateAlertHandlerSlack(alertHandlerSlack AlertHandler, accountId string, handlerId string) (*AlertHandler, error) {
	rb, err := json.Marshal(alertHandlerSlack)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/account/%s/alerting/handlers/slack/%s/", c.HostURL, accountId, handlerId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newResource := AlertHandler{}
	err = json.Unmarshal(body, &newResource)
	if err != nil {
		return nil, err
	}

	return &newResource, nil
}

func (c *Client) DeleteAlertHandlerSlack(AlertHandlerSlackUuid string, accountId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/alerting/handlers/slack/%s/", c.HostURL, accountId, AlertHandlerSlackUuid), nil)
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
