package main

import (
	"fmt"
	"henzworld/handlers"
	"henzworld/shared"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Config struct {
	GHToken       string
	LetterboxdURL string
	StatusCafeURL string
	GoodreadsURL  string
	Port          string
	CacheEnabled  bool
}

func loadConfig() (*Config, error) {
	ghToken, err := getEnvOrError("GITHUB_TOKEN")
	if err != nil {
		return nil, err
	}

	letterboxdURL, err := getEnvOrError("LETTERBOXD_URL")
	if err != nil {
		return nil, err
	}

	statusCafeURL, err := getEnvOrError("STATUS_CAFE_URL")
	if err != nil {
		return nil, err
	}

	goodreadsURL, err := getEnvOrError("GOODREADS_URL")
	if err != nil {
		return nil, err
	}

	cacheEnabledValue := os.Getenv("CACHE_ENABLED")
	cacheEnabled, err := strconv.ParseBool(cacheEnabledValue)
	if err != nil {
		cacheEnabled = true
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		GHToken:       *ghToken,
		LetterboxdURL: *letterboxdURL,
		StatusCafeURL: *statusCafeURL,
		GoodreadsURL:  *goodreadsURL,
		Port:          port,
		CacheEnabled:  cacheEnabled,
	}, nil
}

func getEnvOrError(key string) (*string, error) {
	value := os.Getenv(key)
	if value == "" {
		return nil, fmt.Errorf("%s env variable not set", key)
	}

	return &value, nil
}

func setupClients(config *Config) *handlers.Clients {
	httpClient := &http.Client{
		// TODO timeout config
		Timeout: 5 * time.Second,
	}

	var gitHubHandler handlers.GitHubHandler
	var statusCafeHandler handlers.StatusCafeHandler
	var goodreadsHandler handlers.GoodreadsHandler
	var letterboxdHandler handlers.LetterboxdHandler

	gitHubClient := shared.NewGitHubClient(httpClient, config.GHToken)
	statusCafeClient := shared.NewStatusClient(httpClient, config.StatusCafeURL)
	goodreadsClient := shared.NewRssClient(httpClient, config.GoodreadsURL)
	letterboxdClient := shared.NewRssClient(httpClient, config.LetterboxdURL)

	if config.CacheEnabled {
		cache := shared.NewMemoryCache()

		gitHubHandler = &handlers.CachedGitHubHandler{
			GitHubAPIClient: gitHubClient,
			Cache:           cache,
		}
		statusCafeHandler = &handlers.CachedStatusCafeHandler{
			StatusCafeClient: statusCafeClient,
			Cache:            cache,
		}
		goodreadsHandler = &handlers.CachedGoodreadsHandler{
			Client: goodreadsClient,
			Cache:  cache,
		}

		letterboxdHandler = &handlers.CachedLetterboxdHandler{
			Client: letterboxdClient,
			Cache:  cache,
		}
	} else {
		gitHubHandler = &handlers.DefaultGitHubHandler{
			GitHubAPIClient: gitHubClient,
		}
		statusCafeHandler = &handlers.DefaultStatusCafeHandler{
			StatusCafeClient: statusCafeClient,
		}
		goodreadsHandler = &handlers.DefaultGoodreadsHandler{
			Client: goodreadsClient,
		}

		letterboxdHandler = &handlers.DefaultLetterboxdHandler{
			Client: letterboxdClient,
		}
	}

	clients := handlers.Clients{
		GitHub:     gitHubHandler,
		Letterboxd: letterboxdHandler,
		StatusCafe: statusCafeHandler,
		Goodreads:  goodreadsHandler,
	}
	return &clients
}

func main() {
	templates, err := template.ParseGlob("template/*.html")
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}

	config, err := loadConfig()
	if err != nil {
		log.Fatal("Could not load config:", err)
	}

	clients := setupClients(config)

	homeHandler := handlers.NewHomeHandler(clients, templates)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/{$}", homeHandler)

	log.Printf("Server starting on :%s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
