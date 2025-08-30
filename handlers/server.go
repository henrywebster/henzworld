package handlers

import (
	"henzworld/shared"
	"html/template"
	"log"
	"net/http"
)

func NewHomeHandler(clients *Clients, templates *template.Template) http.HandlerFunc {
	//gitHubAPICommitFetcher := GitHubAPIFetcher{
	//	GitHubAPIClient: clients.GitHub,
	//}
	//
	//

	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title            string
			Message          string
			Commits          []shared.Commit
			RecentlyWatched  []shared.Movie
			Status           *shared.Status
			CurrentlyReading []shared.Book
		}{
			Title:   "home",
			Message: "henz world",
		}

		// Commits

		commits, err := clients.GitHub.GetCommits()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data.Commits = commits

		// Movies

		feed, err := clients.Letterboxd.GetRecentlyWatched()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data.RecentlyWatched = feed

		// Status

		status, err := clients.StatusCafe.GetStatus()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data.Status = status

		// Currently reading

		books, err := clients.Goodreads.GetCurrentlyReading()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data.CurrentlyReading = books

		// Template and write out

		if err := templates.ExecuteTemplate(w, "home.html", data); err != nil {
			// log the detailed error and return with generic 500
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
