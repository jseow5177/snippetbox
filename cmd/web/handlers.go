package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't, use
	// the http.NotFound() function to send a 404 response to the client.
	// Importantly, we then return from the handler. If we don't return the handler
	// would keep executing and also write the "Hello from SnippetBox" message.
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// Initialize a slice containing the paths to the two files. Note that the
	// home.page.html file must be the *first* file in the slice
	files := []string{
		"./ui/html/home.page.html", // Home page
		"./ui/html/footer.partial.html", // Footer partial
		"./ui/html/base.layout.html", // Base template
	}

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there is an error, we log the detailed error message and use
	// the http.Error() function to send a generic 500 Internal Server Error response
	// to the user. Note that we can pass the slice of file paths as a variadic parameter.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// We then use the Execute() method on the template set to write the template content
	// as the response body. The last parameter to Execute() represents any dynamic data
	// that we want to pass in, which for now we'll set as nil.
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the vakue of id parameter from query string.
	// Try to convert it to an integer using the strconv.Atoi() function.
	// If it can't be converted to an integer, or the value if less than 1, return a
	// 404 page not found response
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not. Note that
	// http.MethodPost is a constant equal to the string "POST"
	if r.Method != http.MethodPost {
		// Add an "Allow: POST" header.
		// Must be before w.WriteHeader and w.Write
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet..."))
}