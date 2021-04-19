package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/jseow5177/snippetbox/pkg/forms"
	"github.com/jseow5177/snippetbox/pkg/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// This is required because Go's html/template package allows us to pass in one
// - and only one - item of dynamic data when rendering a template.
type templateData struct {
	CurrentYear int
	Flash string // Flash message on successful POST
	Form *forms.Form
	Snippet *models.Snippet
	Snippets []*models.Snippet
}

// formatDate() is a custom template function that returns a nicely formatted
// string representation of a time.Time object
func formatDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable. This is essentially
// a string-keyed map which acts as a lookup between the names of our custom template functions
// and the functions themselves.
// Each function can have multiple parameters, but they must have either a single return value, or two
// return values of which the second has type error.
var functions = template.FuncMap{
	"formatDate": formatDate,
}

// A map that acts as a template cache
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.html'. This gives us a slice of all the 
	// 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	// Loop through the pages
	for _, page := range pages {
		// Extract the file name (like 'home.page.html') from the full file path
		// and assign it to the name variable.
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. Hence, we need to use template.New() to create
		// an empty template set, use the Funcs() method to register the template.FuncMap,
		// and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'layout' templates to the template set
		// ParseGlob is similar to ParseFiles except that it supports widcard matching
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'partial' templates to the
		// template set.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache. Use the name of the page
		// as the key
		cache[name] = ts
	}

	return cache, nil
}