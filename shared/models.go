package shared

import (
	"fmt"
	"html"
	"sort"
	"time"

	"github.com/mmcdole/gofeed"
)

type Commit struct {
	Message  string
	URL      string
	RepoName string
	RepoURL  string
	Date     time.Time
}

func (r *GitHubAPIResponse) ToCommits() ([]Commit, error) {
	var commits []Commit

	if len(r.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", r.Errors)
	}

	for _, repo := range r.Data.Viewer.Repositories.Nodes {
		if repo.DefaultBranchRef.Target.History.Nodes == nil {
			// no commits
			continue
		}

		for _, node := range repo.DefaultBranchRef.Target.History.Nodes {
			commitDate, err := time.Parse(time.RFC3339, node.CommittedDate)
			if err != nil {
				return nil, fmt.Errorf("invalid commit date %s: %w", node.CommittedDate, err)
			}

			commit := Commit{
				Message:  node.MessageHeadline,
				URL:      node.CommitURL,
				RepoURL:  repo.URL,
				RepoName: repo.Name,
				Date:     commitDate,
			}

			commits = append(commits, commit)
		}
	}

	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.After(commits[j].Date)
	})

	return commits, nil
}

// RSS

type Movie struct {
	Title       string
	Year        string
	WatchedDate string
	URL         string
}

func GetLetterboxdWatched(items []*gofeed.Item) []Movie {
	var movies []Movie

	for i := range items {
		// formattedDate, _ := formatTime("2006-01-02", items[i].Extensions["letterboxd"]["watchedDate"][0].Value)

		movie := Movie{
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

type Book struct {
	Title      string
	AuthorName string
	URL        string
}

func GetGoodreadsCurrentlyReading(items []*gofeed.Item) []Book {
	var books []Book

	for i := range items {
		// TODO format date
		book := Book{
			Title:      items[i].Title,
			AuthorName: items[i].Custom["author_name"],
			URL:        items[i].Link,
		}
		books = append(books, book)
	}
	return books
}

// Status

type Status struct {
	Content string
	Face    string
	TimeAgo string
}

func (r *StatusResponse) GetStatus() *Status {
	escapedContent := html.UnescapeString(r.Content)

	status := Status{
		Content: escapedContent,
		Face:    r.Face,
		TimeAgo: r.TimeAgo,
	}

	return &status
}
