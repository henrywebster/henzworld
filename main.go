package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var templates *template.Template

type Config struct {
	GHToken       string
	LetterboxdURL string
	StatusCafeURL string
	GoodreadsURL  string
	Port          string
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
	}, nil
}

func getEnvOrError(key string) (*string, error) {
	value := os.Getenv(key)
	if value == "" {
		return nil, fmt.Errorf("%s env variable not set", key)
	}

	return &value, nil
}

type Clients struct {
	GitHub     *GitHubClient
	Letterboxd *RssClient
	StatusCafe *StatusClient
	Goodreads  *RssClient
}

func main() {
	httpClient := &http.Client{
		// TODO timeout config
		Timeout: 5 * time.Second,
	}

	var err error
	templates, err = template.ParseGlob("template/*.html")
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}

	config, err := loadConfig()
	if err != nil {
		log.Fatal("Could not load config:", err)
	}

	clients := Clients{
		GitHub:     NewGitHubClient(httpClient, config.GHToken),
		Letterboxd: NewRssClient(httpClient, config.LetterboxdURL),
		StatusCafe: NewStatusClient(httpClient, config.StatusCafeURL),
		Goodreads:  NewRssClient(httpClient, config.GoodreadsURL),
	}

	homeHandler := newHomeHandler(&clients)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/{$}", homeHandler)

	log.Printf("Server starting on :%s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
