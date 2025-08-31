// Package model: Domain models
package model

import (
	"time"
)

type Commit struct {
	Message  string
	URL      string
	RepoName string
	RepoURL  string
	Date     time.Time
}

type Movie struct {
	Title       string
	Year        string
	WatchedDate time.Time
	URL         string
}

type Book struct {
	Title      string
	AuthorName string
	URL        string
}

type Status struct {
	Content string
	Face    string
	TimeAgo string
}

type Post struct {
	Title       string
	CreatedAt   time.Time
	Description string
}
