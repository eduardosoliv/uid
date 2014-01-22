package main

import (
	"encoding/json"
	"errors"
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
var typesJson = flag.String(
	"types",
	`{"test": "uid_test", "sample": "uid_sample"}`,
	`The types as JSON eg: {"test": "uid_test", "sample": "uid_sample"}`,
)
var maxAmount = flag.Int("max-amount", 100, "The max amount of number of ids")
var types map[string]string

func handler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path[0:]
	method := r.Method

	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic: ", err)
			errorHandler(w, r, http.StatusInternalServerError, nil)
		}
	}()

	// check URI and method
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

	var err error

	// get amount
	var amount int
	if amount, err = getAmount(r); err != nil {
		errorHandlerFromError(w, r, http.StatusBadRequest, err)
		return
	}

	// get type and table
	var idsType, table string
	if idsType, table, err = getTypeAndTable(r); err != nil {
		errorHandlerFromError(w, r, http.StatusBadRequest, err)
		return
	}

	// get format and content type
	format, contentType := getFormat(r)

	fmt.Fprintf(w, "%d %s %s %s %s", amount, idsType, table, format, contentType)
}

func getAmount(r *http.Request) (int, error) {
	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil || amount < 1 {
		amount = 1
	}

	if amount <= *maxAmount {
		return amount, nil
	}

	return amount, errors.New(
		fmt.Sprintf("Amount %d exceedes maximum amount %d", amount, *maxAmount),
	)
}

func getTypeAndTable(r *http.Request) (string, string, error) {
	idsType := r.URL.Path[utf8.RuneCountInString(*uri):]
	table, found := types[idsType]

	var err error
	if !found {
		err = errors.New(
			fmt.Sprintf("Type %s not recognized", idsType),
		)
	}

	return idsType, table, err
}

func storeTypes() {
	js := json.NewDecoder(strings.NewReader(strings.TrimSpace(*typesJson)))

	if err := js.Decode(&types); err != nil {
		log.Fatalf("Failed decode types json: %s", err)
	}
}

func main() {
	port := flag.Int("port", 8080, "The port to listen")
	flag.Parse()
	storeTypes()

	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
