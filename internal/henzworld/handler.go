// Package henzworld: core server logic
package henzworld

import (
	"henzworld/internal/database"
	"henzworld/internal/model"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
)

func NewHomeHandler(clients *Clients, templates *template.Template, navEnabled bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title            string
			Message          string
			Commits          []model.Commit
			RecentlyWatched  []model.Movie
			Status           *model.Status
			CurrentlyReading []model.Book
			Page             string
			NavEnabled       bool
		}{
			Title:      "home",
			Message:    "henz world",
			Page:       "Home",
			NavEnabled: navEnabled,
		}

		// Commits

		commits, err := clients.GitHub.GetCommits()
		if err != nil {
			log.Print(err)
		} else {
			data.Commits = commits
		}

		// Movies

		feed, err := clients.Letterboxd.GetRecentlyWatched()
		if err != nil {
			log.Print(err)
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
		} else {
			data.CurrentlyReading = books
		}

		// Template and write out

		if err := templates.ExecuteTemplate(w, "layout", data); err != nil {
			// log the detailed error and return with generic 500
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func NewBlogHandler(db *database.DB, blogTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := db.GetPosts()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		data := struct {
			Posts      []model.Post
			Page       string
			NavEnabled bool
		}{
			Posts:      posts,
			Page:       "Blog",
			NavEnabled: true,
		}

		if err := blogTemplate.ExecuteTemplate(w, "layout", data); err != nil {
			// log the detailed error and return with generic 500
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func NewBlogPostHandler(db *database.DB, postTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		post, err := db.GetPost(slug)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := struct {
			Post       *model.Post
			Page       string
			NavEnabled bool
		}{
			Post:       post,
			Page:       "Blog",
			NavEnabled: true,
		}

		postTemplate.ExecuteTemplate(w, "layout", data)
	}
}

func NewBlogFeedHandler(db *database.DB, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/atom+xml")

		now := time.Now()
		feed := &feeds.Feed{
			Title:       "henz.world blog",
			Link:        &feeds.Link{Href: baseURL + "/blog"},
			Description: "musings from the blog of Henry J. Webster",
			Author:      &feeds.Author{Name: "Henry J. Webster", Email: "henz.world.qv1ok@dralias.com"},
			Created:     now,
		}

		posts, err := db.GetPosts()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		items := postsToFeed(posts, baseURL)

		feed.Items = items

		atom, err := feed.ToAtom()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(atom))
	}
}
