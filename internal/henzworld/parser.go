package henzworld

import (
	"henzworld/internal/model"

	"github.com/gorilla/feeds"
)

func postsToFeed(posts []model.Post, baseURL string) []*feeds.Item {
	var items []*feeds.Item

	for i := range posts {
		item := &feeds.Item{
			Title:       posts[i].Title,
			Link:        &feeds.Link{Href: baseURL + "/blog/post/" + posts[i].Slug},
			Description: posts[i].Description,
			Created:     posts[i].CreatedAt,
		}
		items = append(items, item)
	}
	return items
}
