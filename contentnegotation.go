package main

import (
	"net/http"
)

// A really incomplete implementation of content negotation but good enough
// for now
func getFormat(r *http.Request) (string, string) {
	formats := map[string]string{
		// html
		"text/html":             "html",
		"application/xhtml+xml": "html",
		// txt
		"text/plain": "txt",
		// js
		"application/javascript":   "js",
		"application/x-javascript": "js",
		"text/javascript":          "js",
		// css
		"text/css": "css",
		// json
		"application/json":   "json",
		"application/x-json": "json",
		// xml
		"text/xml":          "xml",
		"application/xml":   "xml",
		"application/x-xml": "xml",
		// rdf
		"application/rdf+xml": "rdf",
		// atom
		"application/atom+xml": "atom",
		// rss
		"application/rss+xml": "rss",
	}

	fallbackFormat := "json"
	fallbackContentType := "application/json"

	supportedFormats := map[string]bool{
		"json": true,
		"xml":  true,
	}

	contentTypes := map[string]string{
		"json": "application/json",
		"xml":  "text/xml",
	}

	// choose format
	accept := r.Header.Get("Accept")
	format, ok := formats[accept]
	if !ok {
		format = "json"
	}

	// check if is supported
	supported, ok := supportedFormats[format]
	if !supported || !ok {
		format = fallbackFormat
	}

	// get content type
	contentType, ok := contentTypes[format]
	if !ok {
		return fallbackFormat, fallbackContentType
	}

	return format, contentType
}
