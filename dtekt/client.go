package dtekt

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Client -

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

type RunConfiguration struct {
	RunConfiguration Runs `json:"runConfiguration"`
}

type Runs struct {
	Runs []Run `json:"runs,omitempty"`
}

type Run struct {
	UUID     string `json:"uuid"`
	Url      string `json:"url"`
	Schedule int    `json:"schedule"`
	Location string `json:"location"`
	Count    int    `json:"count"`
}

type Account struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UptimeMonitor struct {
	UUID            string              `json:"uuid"`
	Url             string              `json:"url"`
	MatchString     string              `json:"match_string"`
	SuccessCode     int                 `json:"success_code"`
	Every           string              `json:"every"`
	Headers         map[string]string   `json:"headers"`
	AlertDefinition []AlertDefinition   `json:"alert_definitions"`
	LocationConfig  map[string][]string `json:"location_config"`
}

type PerformanceMonitor struct {
	UUID            string              `json:"uuid"`
	Url             string              `json:"url"`
	Every           string              `json:"every"`
	AlertDefinition []AlertDefinition   `json:"alert_definitions"`
	LocationConfig  map[string][]string `json:"location_config"`
}

type AlertHandler struct {
	UUID    string                 `json:"uuid"`
	Kind    string                 `json:"kind"`
	Options map[string]interface{} `json:"options"`
	Name    string                 `json:"name"`
}

type AlertDefinition struct {
	UUID      string        `json:"uuid"`
	Warn      float64       `json:"warn"`
	Crit      float64       `json:"crit"`
	Metric    string        `json:"metric"`
	Handlers  []interface{} `json:"handlers"`
	MonitorId string        `json:"monitor_id"`
	Window    string        `json:"window"`
}

// NewClient -
func NewClient(host, api_token *string) (*Client, error) {
	apiurl, exists := os.LookupEnv("DTEKT_API_URL")

	var HostURL string

	if exists {
		HostURL = apiurl
	} else {
		HostURL = "https://siteperfapi.kscloud.pl"
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: 600 * time.Second},
		// Default DTEKT.IO API URL
		HostURL: HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	c.Token = *api_token

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.Token))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	statusOK := res.StatusCode >= 200 && res.StatusCode < 300

	if !statusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
