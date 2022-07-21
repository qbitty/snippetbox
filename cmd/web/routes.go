package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/qbitty/snippetbox/pkg/config"
)

func routes(app *config.Application) http.Handler {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(home(app)))
	mux.Get("/snippet/create", http.HandlerFunc(createSnippetForm(app)))
	mux.Post("/snippet/create", http.HandlerFunc(createSnippet(app)))
	mux.Get("/snippet/:id", http.HandlerFunc(showSnippet(app)))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return recoverPanic(app, logRequest(app, secureHeaders(mux)))

	// standardMiddleware := alice.New(secureHeaders)
	// return standardMiddleware.Then(mux)
}
