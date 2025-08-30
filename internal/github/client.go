// Package github: GitHub API
package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"henzworld/internal"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	token      string
}

func NewClient(httpClient *http.Client, token string) *Client {
	return &Client{
		httpClient: httpClient,
		token:      token,
	}
}

type Request struct {
	Query string `json:"query"`
}

type Response struct {
	Data struct {
		Viewer struct {
			Repositories struct {
				Nodes []struct {
					Name             string `json:"name"`
					URL              string `json:"url"`
					UpdatedAt        string `json:"updatedAt"`
					DefaultBranchRef struct {
						Target struct {
							History struct {
								Nodes []struct {
									MessageHeadline string `json:"messageHeadline"`
									CommittedDate   string `json:"committedDate"`
									CommitURL       string `json:"commitUrl"`
								} `json:"nodes"`
							} `json:"history"`
						} `json:"target"`
					} `json:"defaultBranchRef"`
				} `json:"nodes"`
			} `json:"repositories"`
		} `json:"viewer"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

func (c *Client) GetPublicRepoCommits() (*Response, error) {
	defer internal.TimeFunction("github_public_repo_commits")()
	reqBody := Request{
		Query: PublicReposCommitsQuery,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var ghResp Response
	if err := json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return nil, err
	}

	return &ghResp, nil
}
