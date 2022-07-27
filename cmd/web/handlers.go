package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/qbitty/snippetbox/pkg/config"
	"github.com/qbitty/snippetbox/pkg/forms"
	"github.com/qbitty/snippetbox/pkg/models"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results, err := app.Snippets.Latest()
		if err != nil {
			serverError(app, w, err)
			return
		}

		render(app, w, r, "home.page.tmpl", &templateData{
			Snippets: results,
		})
	}
}

// Add a showSnippet handler function.
func showSnippet(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get(":id"))
		if err != nil || id < 1 {
			notFound(app, w)
			return
		}

		snippet, err := app.Snippets.Get(id)
		if err == models.ErrNoRecord {
			notFound(app, w)
			return
		} else if err != nil {
			serverError(app, w, err)
			return
		}

		// Use the new render helper.
		render(app, w, r, "show.page.tmpl", &templateData{
			Snippet: snippet,
		})
	}
}

// Add a createSnippet handler function.
func createSnippet(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			clientError(app, w, http.StatusBadRequest)
			return
		}

		form := forms.New(r.PostForm)
		form.Required("title", "content", "expires")
		form.MaxLength("title", 100)
		form.PermittedValues("expires", "365", "7", "1")

		if !form.Valid() {
			render(app, w, r, "create.page.tmpl", &templateData{Form: form})
			return
		}

		id, err := app.Snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
		if err != nil {
			serverError(app, w, err)
			return
		}

		app.Session.Put(r, "flash", "Snippet successfully created!")

		http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
	}
}

func createSnippetForm(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(app, w, r, "create.page.tmpl", &templateData{
			// Pass a new empty forms.Form object to the template.
			Form: forms.New(nil),
		})
	}
}
