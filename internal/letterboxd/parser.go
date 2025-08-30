package letterboxd

import (
	"henzworld/internal/model"

	"github.com/mmcdole/gofeed"
)

func GetLetterboxdWatched(items []*gofeed.Item) []model.Movie {
	var movies []model.Movie

	for i := range items {
		// formattedDate, _ := formatTime("2006-01-02", items[i].Extensions["letterboxd"]["watchedDate"][0].Value)

		movie := model.Movie{
			Title:       items[i].Extensions["letterboxd"]["filmTitle"][0].Value,
			Year:        items[i].Extensions["letterboxd"]["filmYear"][0].Value,
			WatchedDate: items[i].Extensions["letterboxd"]["watchedDate"][0].Value,
			URL:         items[i].Link,
		}
		movies = append(movies, movie)
	}

	// TODO sort by WatchedDate

	if len(movies) > 5 {
		movies = movies[:5]
	}
	return movies
}
