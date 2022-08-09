package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func routes(app *application) http.Handler {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.Session.Enable, noSurf, app.authenticate)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(home(app)))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(createSnippetForm(app)))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(createSnippet(app)))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(showSnippet(app)))

	// Add the five new routes.
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(signupUserForm(app)))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(signupUser(app)))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(loginUserForm(app)))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(loginUser(app)))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(logoutUser(app)))

	mux.Get("/ping", http.HandlerFunc(ping))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
