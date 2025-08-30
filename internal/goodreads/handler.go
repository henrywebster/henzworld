package goodreads

import (
	"henzworld/internal"
	"henzworld/internal/model"
	"time"
)

type Handler interface {
	GetCurrentlyReading() ([]model.Book, error)
}

type DefaultHandler struct {
	Client *internal.RssClient
}

type CachedHandler struct {
	Client *internal.RssClient
	Cache  *internal.MemoryCache
}

func (handler *DefaultHandler) GetCurrentlyReading() ([]model.Book, error) {
	goodreadsFeed, err := handler.Client.GetRssFeed()
	if err != nil {
		return nil, err
	}

	books := GetGoodreadsCurrentlyReading(goodreadsFeed)

	return books, nil
}

func (handler *CachedHandler) GetCurrentlyReading() ([]model.Book, error) {
	cacheKey := "currently_reading"
	if cached, found := handler.Cache.Get(cacheKey); found {
		return cached.([]model.Book), nil
	}

	goodreadsFeed, err := handler.Client.GetRssFeed()
	if err != nil {
		return nil, err
	}

	books := GetGoodreadsCurrentlyReading(goodreadsFeed)

	handler.Cache.Set(cacheKey, books, 5*time.Minute)
	return books, nil
}
