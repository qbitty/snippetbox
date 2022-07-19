package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/qbitty/snippetbox/pkg/config"
	"github.com/qbitty/snippetbox/pkg/models"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			notFound(app, w)
			return
		}

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
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
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
		if r.Method != "POST" {
			w.Header().Set("Allow", "POST")
			// w.WriteHeader(405)
			// w.Write([]byte("Method Not Allowed\n"))
			clientError(app, w, http.StatusMethodNotAllowed)
			return
		}
		title := "O snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi"
		expires := "7"

		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			serverError(app, w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
	}
}
