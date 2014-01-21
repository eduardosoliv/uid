package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

// note, that variables are pointers
var uri = flag.String("uri", "/api/ids/", "The URI")

func handler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path[0:]
	method := r.Method

	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic: ", err)
			errorHandler(w, r, http.StatusInternalServerError, nil)
		}
	}()

	if strings.Index(urlPath, *uri) != 0 || method != "GET" {
		errorHandler(
			w,
			r,
			http.StatusNotFound,
			[]string{
				fmt.Sprintf(
					"No route found for \"%s %s\"",
					method,
					urlPath,
				),
			},
		)
		return
	}

	amount := getAmount(r)
	fmt.Fprintf(w, "%s %d", getType(r), amount)

	format, contentType := getFormat(r)

	fmt.Fprintf(w, "%s %s %s", urlPath, format, contentType)
}

func getType(r *http.Request) string {
	return r.URL.Path[utf8.RuneCountInString(*uri):]
}

func getAmount(r *http.Request) int {
	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil || amount < 1 {
		amount = 1
	}

	return amount
}

func main() {
	port := flag.Int("port", 8080, "The port to listen")
	flag.Parse()
	println(*uri)
	println(*port)
	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
