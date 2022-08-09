package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/qbitty/snippetbox/pkg/forms"
	"github.com/qbitty/snippetbox/pkg/models"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(app *application) http.HandlerFunc {
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
func showSnippet(app *application) http.HandlerFunc {
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
func createSnippet(app *application) http.HandlerFunc {
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

func createSnippetForm(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(app, w, r, "create.page.tmpl", &templateData{
			// Pass a new empty forms.Form object to the template.
			Form: forms.New(nil),
		})
	}
}

func signupUserForm(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(app, w, r, "signup.page.tmpl", &templateData{
			Form: forms.New(nil),
		})
	}
}

func signupUser(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the form data.
		err := r.ParseForm()
		if err != nil {
			clientError(app, w, http.StatusBadRequest)
			return
		}
		// Validate the form contents using the form helper we made earlier.
		form := forms.New(r.PostForm)
		form.Required("name", "email", "password")
		form.MaxLength("name", 255)
		form.MaxLength("email", 255)
		form.MatchesPattern("email", forms.EmailRX)
		form.MinLength("password", 10)
		// If there are any errors, redisplay the signup form.
		if !form.Valid() {
			render(app, w, r, "signup.page.tmpl", &templateData{Form: form})
			return
		}

		err = app.Users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
		if err != nil {
			if errors.Is(err, models.ErrDuplicateEmail) {
				form.Errors.Add("email", "Address is already in use")
				render(app, w, r, "signup.page.tmpl", &templateData{Form: form})
			} else {
				serverError(app, w, err)
			}
			return
		}

		app.Session.Put(r, "flash", "Your signup was successful. Please log in.")
		// And redirect the user to the login page.
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
}

func loginUserForm(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(app, w, r, "login.page.tmpl", &templateData{
			Form: forms.New(nil),
		})
	}
}

func loginUser(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			clientError(app, w, http.StatusBadRequest)
			return
		}
		// Check whether the credentials are valid. If they're not, add a generic error
		// message to the form failures map and re-display the login page.
		form := forms.New(r.PostForm)
		id, err := app.Users.Authenticate(form.Get("email"), form.Get("password"))
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				form.Errors.Add("generic", "Email or Password is incorrect")
				render(app, w, r, "login.page.tmpl", &templateData{Form: form})
			} else {
				serverError(app, w, err)
			}
			return
		}
		// Add the ID of the current user to the session, so that they are now 'logged
		// in'.
		app.Session.Put(r, "authenticatedUserID", id)
		// Redirect the user to the create snippet page.
		http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
	}
}

func logoutUser(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 'logged out'.
		app.Session.Remove(r, "authenticatedUserID")
		// Add a flash message to the session to confirm to the user that they've been
		// logged out.
		app.Session.Put(r, "flash", "You've been logged out successfully!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
