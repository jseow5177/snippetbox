package main

import (
	"net/http"
	"path/filepath"

	"github.com/justinas/alice"
)

// Implement custom file system
type neuteredFileSystem struct {
	fs http.FileSystem
}

// A custom FileSystem to disable directory listings for requests to /static/
// Will serve index.html in /static/ if present
// https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	// Stat returns a FileInfo describing the named file from the file system
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	// Checks if path is a directory
	if s.IsDir() {
		// filepath is OS aware. In Windows, filepath.Join uses backslash.
		// ToSlash() replaces the separator to a slash "/" character.
		index := filepath.ToSlash(filepath.Join(path, "/index.html"))

		_, err := nfs.fs.Open(index) // Checks if there is an index.html file
		if err != nil {
			closeErr := f.Close() // Close file
			if closeErr != nil {
				return nil, closeErr
			}
			
			// Return error if no index.html instead of showing directory listing
			// error will be transformed into a 404 Not Found by http.FileServer
			return nil, err
		}
	}

	return f, nil
}

func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// A custom file system that disables directory listing
	customFs := neuteredFileSystem {
		fs: http.Dir(app.config.StaticDir),
	}

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(customFs)

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	return standardMiddleware.Then(mux)
}