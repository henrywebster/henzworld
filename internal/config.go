// Package internal: shared utilities
package internal

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	GHToken       string
	LetterboxdURL string
	StatusCafeURL string
	GoodreadsURL  string
	Port          string
	CacheEnabled  bool
}

func LoadConfig() (*Config, error) {
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
