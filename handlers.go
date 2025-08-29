package main

import (
	"log"
	"net/http"
)

/*
* Right now, homeHandler takes care of everything.
* But that will change.
*
 */
func newHomeHandler(clients *Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title            string
			Message          string
			Commits          []Commit
			RecentlyWatched  []Movie
			Status           *Status
			CurrentlyReading []Book
		}{
			Title:   "home",
			Message: "henz world",
		}

		// Commits

		ghResponse, err := clients.GitHub.GetPublicRepoCommits()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		commits, err := ghResponse.ToCommits()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if len(commits) > 5 {
			commits = commits[:5]
		}

		data.Commits = commits

		// Movies

		feed, err := clients.Letterboxd.GetRssFeed()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data.RecentlyWatched = getLetterboxdWatched(feed)

		// Status

		statusResponse, err := clients.StatusCafe.GetStatus()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data.Status = statusResponse.GetStatus()

		// Currently reading

		goodreadsFeed, err := clients.Goodreads.GetRssFeed()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data.CurrentlyReading = getGoodreadsCurrentlyReading(goodreadsFeed)

		// Template and write out

		if err := templates.ExecuteTemplate(w, "home.html", data); err != nil {
			// log the detailed error and return with generic 500
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
