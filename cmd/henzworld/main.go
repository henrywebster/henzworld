package main

import (
	"henzworld/internal"
	"henzworld/internal/henzworld"
	"html/template"
	"log"
	"net/http"
)

func main() {
	templates, err := template.ParseGlob("template/*.html")
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}

	config, err := internal.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config:", err)
	}

	clients := henzworld.SetupClients(config)

	homeHandler := henzworld.NewHomeHandler(clients, templates)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/{$}", homeHandler)

	log.Printf("Starting henzworld on :%s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
