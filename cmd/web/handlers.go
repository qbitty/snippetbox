package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/qbitty/snippetbox/pkg/config"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			notFound(app, w)
			return
		}

		files := []string{
			"./ui/html/home.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/footer.partial.tmpl",
		}

		ts, err := template.ParseFiles(files...)

		if err != nil {
			app.ErrLog.Println(err.Error())
			serverError(app, w, err)
		}

		err = ts.Execute(w, nil)
		if err != nil {
			app.ErrLog.Println(err.Error())
			serverError(app, w, err)
		}
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

		fmt.Fprintf(w, "Display a specific snippet with ID %d...\n", id)
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
		w.Write([]byte("Create a new snippet...\n"))
	}
}
