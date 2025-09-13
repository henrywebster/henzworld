package statuscafe

import (
	"henzworld/internal/model"
)

type Handler interface {
	GetStatus() (*model.Status, error)
}

type DefaultHandler struct {
	Client *Client
}

func (handler *DefaultHandler) GetStatus() (*model.Status, error) {
	statusResponse, err := handler.Client.GetStatus()
	if err != nil {
		return nil, err
	}
	return statusResponse.GetStatus(), nil
}
