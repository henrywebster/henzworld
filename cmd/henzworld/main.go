package main

import (
	"henzworld/internal"
	"henzworld/internal/henzworld"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	config, err := internal.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config:", err)
	}

	templatePattern := filepath.Join(config.TemplateDir, "*.html")
	templates, err := template.ParseGlob(templatePattern)
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
