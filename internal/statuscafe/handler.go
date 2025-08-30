package statuscafe

import (
	"henzworld/internal"
	"henzworld/internal/model"
	"time"
)

type Handler interface {
	GetStatus() (*model.Status, error)
}

type DefaultHandler struct {
	Client *Client
}

type CachedHandler struct {
	Client *Client
	Cache  *internal.MemoryCache
	TTL    time.Duration
}

func (handler *DefaultHandler) GetStatus() (*model.Status, error) {
	statusResponse, err := handler.Client.GetStatus()
	if err != nil {
		return nil, err
	}
	return statusResponse.GetStatus(), nil
}

func (handler *CachedHandler) GetStatus() (*model.Status, error) {
	cacheKey := "status"
	if cached, found := handler.Cache.Get(cacheKey); found {
		return cached.(*model.Status), nil
	}

	statusResponse, err := handler.Client.GetStatus()
	if err != nil {
		return nil, err
	}

	status := statusResponse.GetStatus()

	handler.Cache.Set(cacheKey, status, handler.TTL)
	return status, nil
}
