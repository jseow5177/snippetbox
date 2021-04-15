package main

import (
	"html/template"
	"path/filepath"

	"github.com/jseow5177/snippetbox/pkg/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// This is required because Go's html/template package allows us to pass in one
// - and only one - item of dynamic data when rendering a template.
type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
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

		// Parse the page template file into a template set.
		ts, err := template.ParseFiles(page)
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