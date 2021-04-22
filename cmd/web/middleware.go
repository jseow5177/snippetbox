package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jseow5177/snippetbox/pkg/models"
	"github.com/justinas/nosurf"
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

// A middleware to recover the goroutine handling a request from panic
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

// A middleware to authenticate user by checking if
// 1. User session has a authenticatedUserID
// 2. User exists in database (valid user id)
// 3. User is active
// - When we don't have an authenticated-and-active user, we pass the original and
//   unchanged *http.Request to the next handler in the chain.
// - When we do have an authenticated-and-active user, we create a copy of the request
//   with a contentKeyIsAuthenticated key and true value stored in the request context.
//   We then pass this copy of the *http.Request to the next handler in the chain.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func ( w http.ResponseWriter, r *http.Request) {
		// Check if an authenticatedUserID value exists in the session
		exists := app.session.Exists(r, "authenticatedUserID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		// Fetch the details of the currebt==nt user from the database. If no matching
		// record is found, or the current user has been deactivated, remove (invalid)
		// authenticatedUserID value from their session.
		user, err := app.users.Get(app.session.GetInt(r, "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) || !user.Active {
			app.session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		// Otherwise, the request is coming from an active, authenticated user.
		// We create a new copy of the request, with a true boolean value added to the request context to indicate
		// that the user is authenticated. Then, we call the next handler in the chain *using the new copy of the request*
		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// A middleware to prevent unauthenticated user from entering routes that require authentication
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache
		// (or other intermediary cache)
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// A middleware function to create a customized CSRF cookie with the Secure, Path and HttpOnly flags set
func noSurf(next http.Handler) http.Handler {
	// Construct a new CSRFHandler that calls the specified handler (next) if the CSRF check succeeds.
	csrfHandler := nosurf.New(next)
	// Set the base cookie to use when building a CSRF token cookie.
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true, // So that cookie is not accessible by scripts
		Path: "/", // The URL path that must exist in the requested URL in order to send the Cookie header
		Secure: true, // Sent only through HTTPS
	})

	return csrfHandler
}
