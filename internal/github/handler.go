package github

import (
	"henzworld/internal"
	"henzworld/internal/model"
	"time"
)

type Handler interface {
	GetCommits() ([]model.Commit, error)
}

type DefaultHandler struct {
	Client *Client
}

type CachedHandler struct {
	Client *Client
	Cache  *internal.MemoryCache
	TTL    time.Duration
}

func (handler *DefaultHandler) GetCommits() ([]model.Commit, error) {
	ghResponse, err := handler.Client.GetPublicRepoCommits()
	if err != nil {
		return nil, err
	}

	commits, err := ghResponse.ToCommits()
	if err != nil {
		return nil, err
	}

	if len(commits) > 5 {
		commits = commits[:5]
	}

	return commits, nil
}

func (handler *CachedHandler) GetCommits() ([]model.Commit, error) {
	cacheKey := "github_commits"
	if cached, found := handler.Cache.Get(cacheKey); found {
		return cached.([]model.Commit), nil
	}

	ghResponse, err := handler.Client.GetPublicRepoCommits()
	if err != nil {
		return nil, err
	}

	commits, err := ghResponse.ToCommits()
	if err != nil {
		return nil, err
	}

	if len(commits) > 5 {
		commits = commits[:5]
	}

	handler.Cache.Set(cacheKey, commits, handler.TTL)
	return commits, nil
}
