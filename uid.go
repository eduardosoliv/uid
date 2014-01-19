package main

import (
	"fmt"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.Path[0:]
	var method = r.Method

	// check that the URL is /api/i/ids/{type}?amount=X
	var index = strings.Index(path, "/api/i/ids/")
	if index != 0 || method != "GET" {
		errorHandler(
			w,
			r,
			http.StatusNotFound,
			[]string{fmt.Sprintf(
				"No route found for \"%s %s\"",
				r.Method,
				r.URL.Path[0:],
			)},
		)
		return
	}

	format, contentType := getFormat(r)

	fmt.Fprintf(w, "%s %s %s", path, format, contentType)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
