package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mmcdole/gofeed"
)

// GitHub

type GitHubClient struct {
	httpClient *http.Client
	token      string
}

func NewGitHubClient(httpClient *http.Client, token string) *GitHubClient {
	return &GitHubClient{
		httpClient: httpClient,
		token:      token,
	}
}

type GraphQLRequest struct {
	Query string `json:"query"`
}

type GitHubAPIResponse struct {
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

func (c *GitHubClient) GetPublicRepoCommits() (*GitHubAPIResponse, error) {
	reqBody := GraphQLRequest{
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

	var ghResp GitHubAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return nil, err
	}

	return &ghResp, nil
}

// RSS

type RssClient struct {
	parser *gofeed.Parser
	url    string
}

func NewRssClient(httpClient *http.Client, url string) *RssClient {
	fp := gofeed.NewParser()
	fp.Client = httpClient

	return &RssClient{
		parser: fp,
		url:    url,
	}
}

func (c *RssClient) GetRssFeed() ([]*gofeed.Item, error) {
	feed, err := c.parser.ParseURL(c.url)
	if err != nil {
		return nil, err
	}

	return feed.Items, nil
}

// Status

type StatusResponse struct {
	Content string `json:"content"`
	Face    string `json:"face"`
	TimeAgo string `json:"timeAgo"`
}

type StatusClient struct {
	httpClient *http.Client
	url        string
}

func NewStatusClient(httpClient *http.Client, url string) *StatusClient {
	return &StatusClient{
		httpClient: httpClient,
		url:        url,
	}
}

func (c *StatusClient) GetStatus() (*StatusResponse, error) {
	resp, err := http.Get(c.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var status StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}
