package main

import (
	"bytes"
	"fmt"
	"net/http"
)


func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Retrieve the appropriate template set from the cache based on page name
	// like ('home.page.html'). If no entry exists in the cache with the provided name, 
	// call the serverError helper method
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// Initialize a new buffer
	buf := new(bytes.Buffer)

	// Write the template into buffer instead of straight into the http.ResponseWriter.
	// This serves as a 'trial render' to catch runtime errors in template rendering.
	// If there is an error, we call the serverError helper method and return.
	err := ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write the contents of the buffer to http.ResponseWriter.
	buf.WriteTo(w)
}