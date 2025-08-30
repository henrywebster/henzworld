package handlers

import (
	"henzworld/shared"
	"time"
)

type GitHubHandler interface {
	GetCommits() ([]shared.Commit, error)
}

type DefaultGitHubHandler struct {
	GitHubAPIClient *shared.GitHubClient
}

type CachedGitHubHandler struct {
	GitHubAPIClient *shared.GitHubClient
	Cache           *shared.MemoryCache
}

func (handler *DefaultGitHubHandler) GetCommits() ([]shared.Commit, error) {
	ghResponse, err := handler.GitHubAPIClient.GetPublicRepoCommits()
	if err != nil {
		return nil, err
	}

	commits, err := ghResponse.ToCommits()
	if err != nil {
		return nil, err
	}

	if len(commits) > 5 {
		commits = commits[:5]
	}

	return commits, nil
}

func (handler *CachedGitHubHandler) GetCommits() ([]shared.Commit, error) {
	cacheKey := "github_commits"
	if cached, found := handler.Cache.Get(cacheKey); found {
		return cached.([]shared.Commit), nil
	}

	ghResponse, err := handler.GitHubAPIClient.GetPublicRepoCommits()
	if err != nil {
		return nil, err
	}

	commits, err := ghResponse.ToCommits()
	if err != nil {
		return nil, err
	}

	if len(commits) > 5 {
		commits = commits[:5]
	}

	handler.Cache.Set(cacheKey, commits, 5*time.Minute)
	return commits, nil
}

type StatusCafeHandler interface {
	GetStatus() (*shared.Status, error)
}

type DefaultStatusCafeHandler struct {
	StatusCafeClient *shared.StatusClient
}

type CachedStatusCafeHandler struct {
	StatusCafeClient *shared.StatusClient
	Cache            *shared.MemoryCache
}

func (handler *DefaultStatusCafeHandler) GetStatus() (*shared.Status, error) {
	statusResponse, err := handler.StatusCafeClient.GetStatus()
	if err != nil {
		return nil, err
	}
	return statusResponse.GetStatus(), nil
}

func (fetcher *CachedStatusCafeHandler) GetStatus() (*shared.Status, error) {
	cacheKey := "status"
	if cached, found := fetcher.Cache.Get(cacheKey); found {
		return cached.(*shared.Status), nil
	}

	statusResponse, err := fetcher.StatusCafeClient.GetStatus()
	if err != nil {
		return nil, err
	}

	status := statusResponse.GetStatus()

	fetcher.Cache.Set(cacheKey, status, 5*time.Minute)
	return status, nil
}

type GoodreadsHandler interface {
	GetCurrentlyReading() ([]shared.Book, error)
}

type DefaultGoodreadsHandler struct {
	Client *shared.RssClient
}

type CachedGoodreadsHandler struct {
	Client *shared.RssClient
	Cache  *shared.MemoryCache
}

func (handler *DefaultGoodreadsHandler) GetCurrentlyReading() ([]shared.Book, error) {
	goodreadsFeed, err := handler.Client.GetRssFeed()
	if err != nil {
		return nil, err
	}

	books := shared.GetGoodreadsCurrentlyReading(goodreadsFeed)

	return books, nil
}

func (handler *CachedGoodreadsHandler) GetCurrentlyReading() ([]shared.Book, error) {
	cacheKey := "currently_reading"
	if cached, found := handler.Cache.Get(cacheKey); found {
		return cached.([]shared.Book), nil
	}

	goodreadsFeed, err := handler.Client.GetRssFeed()
	if err != nil {
		return nil, err
	}

	books := shared.GetGoodreadsCurrentlyReading(goodreadsFeed)

	handler.Cache.Set(cacheKey, books, 5*time.Minute)
	return books, nil
}

type LetterboxdHandler interface {
	GetRecentlyWatched() ([]shared.Movie, error)
}

type DefaultLetterboxdHandler struct {
	Client *shared.RssClient
}

type CachedLetterboxdHandler struct {
	Client *shared.RssClient
	Cache  *shared.MemoryCache
}

func (handler *DefaultLetterboxdHandler) GetRecentlyWatched() ([]shared.Movie, error) {
	letterboxdFeed, err := handler.Client.GetRssFeed()
	if err != nil {
		return nil, err
	}

	movies := shared.GetLetterboxdWatched(letterboxdFeed)

	return movies, nil
}

func (handler *CachedLetterboxdHandler) GetRecentlyWatched() ([]shared.Movie, error) {
	cacheKey := "recently_watched"
	if cached, found := handler.Cache.Get(cacheKey); found {
		return cached.([]shared.Movie), nil
	}

	letterboxdFeed, err := handler.Client.GetRssFeed()
	if err != nil {
		return nil, err
	}

	movies := shared.GetLetterboxdWatched(letterboxdFeed)

	handler.Cache.Set(cacheKey, movies, 5*time.Minute)
	return movies, nil
}

type Clients struct {
	GitHub     GitHubHandler
	Letterboxd LetterboxdHandler
	StatusCafe StatusCafeHandler
	Goodreads  GoodreadsHandler
}
