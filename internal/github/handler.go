package github

import (
	"henzworld/internal/model"
)

type Handler interface {
	GetCommits() ([]model.Commit, error)
}

type DefaultHandler struct {
	Client *Client
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
