package dtekt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// func (c *Client) GetPerformanceMonitors(account_id string) (*RunConfiguration, error) {
// 	req, err := http.NewRequest("GET", fmt.Sprintf("%s/performance-monitors/%s/", c.HostURL, account_id), nil)
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

func (c *Client) GetPerformanceMonitor(monitorUuid string, account_id string) (*PerformanceMonitor, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/performance-monitors/%s/", c.HostURL, account_id, monitorUuid), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := PerformanceMonitor{}
	err = json.Unmarshal(body, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) CreatePerformanceMonitor(performanceMonitor PerformanceMonitor, account_id string) (*PerformanceMonitor, error) {
	rb, err := json.Marshal(performanceMonitor)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/performance-monitors/", c.HostURL, account_id), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	monitor := PerformanceMonitor{}
	err = json.Unmarshal(body, &monitor)
	if err != nil {
		return nil, err
	}

	return &monitor, nil
}

func (c *Client) UpdatePerformanceMonitor(performanceMonitor PerformanceMonitor, accountId string, monitorId string) (*PerformanceMonitor, error) {
	rb, err := json.Marshal(performanceMonitor)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/account/%s/performance-monitors/%s/", c.HostURL, accountId, monitorId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	monitor := PerformanceMonitor{}
	err = json.Unmarshal(body, &monitor)
	if err != nil {
		return nil, err
	}

	return &monitor, nil
}

func (c *Client) DeletePerformanceMonitor(PerformanceMonitorUuid string, account_id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/performance-monitors/%s/", c.HostURL, account_id, PerformanceMonitorUuid), nil)
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
