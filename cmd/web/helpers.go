package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	// debug.Stack() function get a stack trace for the current goroutine
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	// Output(calldepth int, s string) writes the output for a logging event.
	// The string s contains the text to print after the prefix specified by the flags of the Logger.
	// Calldepth is the count of the number of frames to skip when computing the file name and line number
	// if Llongfile or Lshortfile is set.
	// For example, a value of 1 will print the details for the caller of Output (which is here)
	app.errorLog.Output(2, trace)

	// http.StatusText converts http error codes to human readable description
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// The notFound helper is a convenience wrapper around clientError which sends a 404 Not Found response
// to the user
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// addDefaultData() injects common dynamic data into our application
// by passing them into an instance of a templateData struct
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CurrentYear = time.Now().Year()

	// Use the PopString() method to retrieve the value for the "flash" key.
	// PopString() also deletes the key and value from the session data, so it 
	// acts like a one-time fetch. If there is no matching key in the session data,
	// this will return an empty string.
	td.Flash = app.session.PopString(r, "flash")

	// Add the authentication status to the template data
	td.IsAuthenticated = app.isAuthenticated(r)

	// Add the CSRF token to the templateData struct
	td.CSRFToken = nosurf.Token(r)

	return td
}

// Return true if the current request is from authenticated user, otherwise return false.
// The isAuthenticated bool, if it exists, will be stored in the request context.
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}