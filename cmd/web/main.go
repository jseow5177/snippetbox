package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

type Config struct {
	Addr string
	StaticDir string
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

// Disable directory listings for requests to /static/ with a custom FileSystem
// https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	// Stat returns a FileInfo describing the named file from the file system
	s, err := f.Stat()
	// Checks if path is a directory
	if s.IsDir() {
		// filepath is OS aware. In Windows, filepath.Join uses backslash.
		// ToSlash() replaces the separator to a slash "/" character.
		index := filepath.ToSlash(filepath.Join(path, "/index.html"))
		fmt.Println(index)
		_, err := nfs.fs.Open(index) // Checks if there is an index.html file
		if err != nil {
			closeErr := f.Close() // Close file
			if closeErr != nil {
				log.Println(closeErr)
				return nil, closeErr
			}
			
			// Return error if no index.html instead of showing directory listing
			// error will be transformed into a 404 Not Found by http.FileServer
			return nil, err
		}
	}

	return f, nil
}

func main() {
	cfg := new(Config)

	// Define a new command-line flag with the name "addr", a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the flag
	// will be stored in the addr variable at runtime
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")

	// We use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr variable.
	// You need to call this *before* you use the addr variable otherwise it will always
	// contain the default value of ":4000". If any errors are encountered during parsing
	// the application will be terminated.
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home) // Works like a "catch-all" route
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// A custom file system that disables directory listing
	customFs := neuteredFileSystem {
		fs: http.Dir(cfg.StaticDir),
	}

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(customFs)

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// The value returned from the flag.String() function is a pointer to the flag value,
	// not the value itself. So we need to dereference the pointer (i.e. prefix it with the * symbol)
	// before using it.
	log.Printf("Starting server on %s", cfg.Addr)
	err := http.ListenAndServe(cfg.Addr, mux)
	log.Fatal(err)
}