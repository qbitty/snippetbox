package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"github.com/qbitty/snippetbox/pkg/config"
)

func routes(app *config.Application) http.Handler {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.

	dynamicMiddleware := alice.New(app.Session.Enable)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(home(app)))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(createSnippetForm(app)))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(createSnippet(app)))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(showSnippet(app)))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return recoverPanic(app, logRequest(app, secureHeaders(mux)))

	// standardMiddleware := alice.New(secureHeaders)
	// return standardMiddleware.Then(mux)
}
