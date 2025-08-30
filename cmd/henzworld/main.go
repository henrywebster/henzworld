package main

import (
	"henzworld/internal"
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

	templatePattern := filepath.Join(config.TemplateDir, "*.html")
	templates, err := template.New("").Funcs(funcMap).ParseGlob(templatePattern)
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}

	clients := henzworld.SetupClients(config)

	homeHandler := henzworld.NewHomeHandler(clients, templates)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.StaticDir))))
	http.HandleFunc("/{$}", homeHandler)

	log.Printf("Starting henzworld on :%s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
