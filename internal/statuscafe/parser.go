package statuscafe

import (
	"henzworld/internal/model"
	"html"
)

func (r *Response) GetStatus() *model.Status {
	escapedContent := html.UnescapeString(r.Content)

	status := model.Status{
		Content: escapedContent,
		Face:    r.Face,
		TimeAgo: r.TimeAgo,
	}

	return &status
}
