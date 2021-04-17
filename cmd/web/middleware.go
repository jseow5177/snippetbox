package main

import (
	"fmt"
	"net/http"
)

// Use the information logger in application dependency to log HTTP requests.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		// RemoteAddr: Network address in the format of "IP:port"
		// Proto: HTTP/1.0
		// Method: HTTP Method
		// URL: 
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)

		next.ServeHTTP(w , r)
	})
}

// A middleware to add headers that prevent XSS attacks 
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent XSS attack
		// Reference: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-XSS-Protection
		w.Header().Set("X-XSS-Protection", "1: mode=block")

		// Prevent Clickjacking
		// Reference: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function which will always be called in the event of a panic
		defer func() {
			// Use the builtin recover function to check if there has been a panic or not
			err := recover()
			if err != nil {
				// Set a "Connection: close" header on the response.
				// Acts as a trigger to make Go's HTTP server automatically close the current connection
				// after a response has been sent.
				// It also informs the user that the connection will be closed.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to returb a 500 Internal Server response
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}