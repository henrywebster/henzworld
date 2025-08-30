package henzworld

import (
	"henzworld/internal"
	"henzworld/internal/github"
	"henzworld/internal/goodreads"
	"henzworld/internal/letterboxd"
	"henzworld/internal/statuscafe"
	"net/http"
	"time"
)

type Clients struct {
	GitHub     github.Handler
	Letterboxd letterboxd.Handler
	StatusCafe statuscafe.Handler
	Goodreads  goodreads.Handler
}

func SetupClients(config *internal.Config) *Clients {
	httpClient := &http.Client{
		// TODO timeout config
		Timeout: 5 * time.Second,
	}

	var gitHubHandler github.Handler
	var statusCafeHandler statuscafe.Handler
	var goodreadsHandler goodreads.Handler
	var letterboxdHandler letterboxd.Handler

	gitHubClient := github.NewClient(httpClient, config.GHToken)
	statusCafeClient := statuscafe.NewClient(httpClient, config.StatusCafeURL)
	goodreadsClient := internal.NewRssClient(httpClient, config.GoodreadsURL)
	letterboxdClient := internal.NewRssClient(httpClient, config.LetterboxdURL)

	if config.CacheEnabled {
		cache := internal.NewMemoryCache()

		gitHubHandler = &github.CachedHandler{
			Client: gitHubClient,
			Cache:  cache,
		}
		goodreadsHandler = &goodreads.CachedHandler{
			Client: goodreadsClient,
			Cache:  cache,
		}
		statusCafeHandler = &statuscafe.CachedHandler{
			Client: statusCafeClient,
			Cache:  cache,
		}
		letterboxdHandler = &letterboxd.CachedHandler{
			Client: letterboxdClient,
			Cache:  cache,
		}
	} else {
		gitHubHandler = &github.DefaultHandler{
			Client: gitHubClient,
		}
		goodreadsHandler = &goodreads.DefaultHandler{
			Client: goodreadsClient,
		}
		statusCafeHandler = &statuscafe.DefaultHandler{
			Client: statusCafeClient,
		}
		letterboxdHandler = &letterboxd.DefaultHandler{
			Client: letterboxdClient,
		}
	}

	clients := Clients{
		GitHub:     gitHubHandler,
		Letterboxd: letterboxdHandler,
		StatusCafe: statusCafeHandler,
		Goodreads:  goodreadsHandler,
	}
	return &clients
}
