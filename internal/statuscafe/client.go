// Package statuscafe: status.cafe API
package statuscafe

import (
	"encoding/json"
	"fmt"
	"henzworld/internal"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	url        string
}

func NewClient(httpClient *http.Client, url string) *Client {
	return &Client{
		httpClient: httpClient,
		url:        url,
	}
}

type Response struct {
	Content string `json:"content"`
	Face    string `json:"face"`
	TimeAgo string `json:"timeAgo"`
}

func (c *Client) GetStatus() (*Response, error) {
	defer internal.TimeFunction("get_status")()
	resp, err := http.Get(c.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var status Response
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}
