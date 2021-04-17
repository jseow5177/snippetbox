package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jseow5177/snippetbox/pkg/forms"
	"github.com/jseow5177/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the render helper method to reder home page
	app.render(w, r, "home.page.html", &templateData{
		Snippets: s,
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Pat doesn't strip thecolon from the named capture key, so we need to get
	// the value of ":id" from the query string instead of "id".
	// Try to convert it to an integer using the strconv.Atoi() function.
	// If it can't be converted to an integer, or the value if less than 1, return a
	// 404 page not found response.
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a specific record
	// based on its ID. If no matching record is found, return a 404 Not Found response.
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.html", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.html", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm() adds any data in POST request bodies to the r.PostForm map.
	// This also works for PUT and PATCH methods.
	// The Content-Type must be application/x-wwww-form-urlencoded.
	// If there is large file data, use ParseMultipartForm where the Content-Type is 
	// multipart/form-data.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new Form struct containing the POSTed data from the form,
	// then use the validation methods to check the content.
	f := forms.New(r.PostForm)
	f.Required("title", "content", "expires")
	f.MaxLength("title", 100)
	f.PermittedValues("expires", "7", "1", "365")

	// If the form isn't valid, redisplay the template passing in the form.Form object as the data.
	if !f.Valid() {
		app.render(w, r, "create.page.html", &templateData{
			Form: f,
		})
		return
	}

	// Because the form data (with type url.Values) has been annonymously embedded
	// in the form.Form struct, we can use the Get() method to retrieve the validated value
	// from a particular form field.
	id, err := app.snippets.Insert(f.Get("title"), f.Get("content"), f.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect user to the page of newly created Snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}