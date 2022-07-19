package main

import (
	"net/http"

	"github.com/qbitty/snippetbox/pkg/config"
)

func routes(app *config.Application) *http.ServeMux {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home(app))
	mux.HandleFunc("/snippet", showSnippet(app))
	mux.HandleFunc("/snippet/create", createSnippet(app))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
