// Package letterboxd: Letterboxd API
package letterboxd

import (
	"henzworld/internal"
	"henzworld/internal/model"
)

type Handler interface {
	GetRecentlyWatched() ([]model.Movie, error)
}

type DefaultHandler struct {
	Client *internal.RssClient
}

func (handler *DefaultHandler) GetRecentlyWatched() ([]model.Movie, error) {
	letterboxdFeed, err := handler.Client.GetRssFeed()
	if err != nil {
		return nil, err
	}

	movies := GetLetterboxdWatched(letterboxdFeed)

	return movies, nil
}
