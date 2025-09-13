package main

import (
	"henzworld/internal"
	"henzworld/internal/database"
	"henzworld/internal/henzworld"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

func main() {
	config, err := internal.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config:", err)
	}

	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("Jan 02, 2006")
		},
		"formatDateTime": func(t time.Time) string {
			return t.Format("Jan 02, 2006 15:04")
		},
		// for the nerds
		"formatHtmlDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"formatHtmlDateTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04")
		},
	}

	homeFiles := []string{
		filepath.Join(config.TemplateDir, "layout.html"),
		filepath.Join(config.TemplateDir, "home.html"),
		filepath.Join(config.TemplateDir, "brief.html"),
		filepath.Join(config.TemplateDir, "commits.html"),
		filepath.Join(config.TemplateDir, "movies.html"),
		filepath.Join(config.TemplateDir, "reading.html"),
		filepath.Join(config.TemplateDir, "status.html"),
	}
	homeTemplate, err := template.New("home").Funcs(funcMap).ParseFiles(homeFiles...)
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}

	clients := henzworld.SetupClients(config)

	navEnabled := config.BlogEnabled

	homeHandler := henzworld.NewHomeHandler(clients, homeTemplate, navEnabled)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.StaticDir))))
	http.HandleFunc("/{$}", homeHandler)

	if config.BlogEnabled {
		db, _ := database.New(config.DatabaseLocalFile)
		blogFiles := []string{
			filepath.Join(config.TemplateDir, "layout.html"),
			filepath.Join(config.TemplateDir, "blog.html"),
		}
		blogTemplate, err := template.New("blog").Funcs(funcMap).ParseFiles(blogFiles...)
		if err != nil {
			log.Fatal("Error loading blog template:", err)
		}

		blogHandler := henzworld.NewBlogHandler(db, blogTemplate)
		http.HandleFunc("/blog/", blogHandler)

		postFiles := []string{
			filepath.Join(config.TemplateDir, "layout.html"),
			filepath.Join(config.TemplateDir, "post.html"),
		}
		postTemplate, err := template.New("post").Funcs(funcMap).ParseFiles(postFiles...)
		if err != nil {
			log.Fatal("Error loading post template:", err)
		}

		postHandler := henzworld.NewBlogPostHandler(db, postTemplate)
		http.HandleFunc("/blog/{slug}/", postHandler)
	}

	log.Printf("Starting henzworld on :%s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
