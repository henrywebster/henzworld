// Package goodreads: Goodreads RSS feed
package goodreads

import (
	"henzworld/internal"
	"henzworld/internal/model"
)

type Handler interface {
	GetCurrentlyReading() ([]model.Book, error)
}

type DefaultHandler struct {
	Client *internal.RssClient
}

func (handler *DefaultHandler) GetCurrentlyReading() ([]model.Book, error) {
	goodreadsFeed, err := handler.Client.GetRssFeed()
	if err != nil {
		return nil, err
	}

	books := GetGoodreadsCurrentlyReading(goodreadsFeed)

	return books, nil
}
