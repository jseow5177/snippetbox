package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
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