package dtekt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) CreateAlertDefinition(alertDefinition AlertDefinition, accountId string) (*AlertDefinition, error) {
	rb, err := json.Marshal(alertDefinition)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/alerting/alertdefinitions", c.HostURL, accountId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newResource := AlertDefinition{}
	err = json.Unmarshal(body, &newResource)
	if err != nil {
		return nil, err
	}

	return &newResource, nil
}

func (c *Client) DeleteAlertDefinition(alertDefinitionUuid string, accountId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/alerting/alertdefinitions/%s/", c.HostURL, accountId, alertDefinitionUuid), nil)
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
