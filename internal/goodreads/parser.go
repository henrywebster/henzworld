package goodreads

import (
	"henzworld/internal/model"

	"github.com/mmcdole/gofeed"
)

func GetGoodreadsCurrentlyReading(items []*gofeed.Item) []model.Book {
	var books []model.Book

	for i := range items {
		// TODO format date
		book := model.Book{
			Title:      items[i].Title,
			AuthorName: items[i].Custom["author_name"],
			URL:        items[i].Link,
		}
		books = append(books, book)
	}
	return books
}
