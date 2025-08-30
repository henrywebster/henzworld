// Package letterboxd: Letterboxd API
package letterboxd

import (
	"henzworld/internal"
	"henzworld/internal/model"
	"time"
)

type Handler interface {
	GetRecentlyWatched() ([]model.Movie, error)
}

type DefaultHandler struct {
	Client *internal.RssClient
}

type CachedHandler struct {
	Client *internal.RssClient
	Cache  *internal.MemoryCache
	TTL    time.Duration
}

func (handler *DefaultHandler) GetRecentlyWatched() ([]model.Movie, error) {
	letterboxdFeed, err := handler.Client.GetRssFeed()
	if err != nil {
		return nil, err
	}

	movies := GetLetterboxdWatched(letterboxdFeed)

	return movies, nil
}

func (handler *CachedHandler) GetRecentlyWatched() ([]model.Movie, error) {
	cacheKey := "recently_watched"
	if cached, found := handler.Cache.Get(cacheKey); found {
		return cached.([]model.Movie), nil
	}

	letterboxdFeed, err := handler.Client.GetRssFeed()
	if err != nil {
		return nil, err
	}

	movies := GetLetterboxdWatched(letterboxdFeed)

	handler.Cache.Set(cacheKey, movies, handler.TTL)
	return movies, nil
}
