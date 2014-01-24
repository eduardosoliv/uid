package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var uri = flag.String("uri", "/api/ids/", "The URI")
var method = flag.String("method", "GET", "The HTTP method")
var typesJson = flag.String(
	"types",
	`{"test": "uid_test", "sample": "uid_sample"}`,
	`The types as JSON eg: {"test": "uid_test", "sample": "uid_sample"}`,
)
var maxAmount = flag.Int("max-amount", 100, "The max amount of number of ids")
var types map[string]string

func handler(w http.ResponseWriter, r *http.Request) {
	var db *sql.DB

	defer func() {
		if err := recover(); err != nil {
			if db != nil {
				db.Close()
			}
			log.Println("Panic: ", err)
			errorHandlerFromError(w, r, http.StatusInternalServerError, err.(error))
		}
	}()

	// check URI and method
	urlPath := r.URL.Path[0:]
	if strings.Index(urlPath, *uri) != 0 || r.Method != *method {
		errorHandler(
			w,
			r,
			http.StatusNotFound,
			[]string{
				fmt.Sprintf(
					"No route found for \"%s %s\"",
					r.Method,
					urlPath,
				),
			},
		)
		return
	}

	// get amount
	var err error
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

	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/eb_unique_test")
	if err != nil {
		panic(err)
	}

	result, err := db.Exec(
		fmt.Sprintf(
			"REPLACE INTO unique_test (stub) VALUES %s",
			getSqlValues(amount),
		),
	)
	if err != nil {
		panic(err)
	}
	firstId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "\n")
	increment := int64(10)
	id := firstId
	for i := 0; i < amount; i++ {
		fmt.Fprintf(w, "%d\n", id)
		id = id + increment
	}

	fmt.Fprintf(w, "%d", id)
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
	idsType := r.FormValue("type")
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

func getSqlValues(amount int) string {
	var buffer bytes.Buffer
	for i := 1; i <= amount; i++ {
		buffer.WriteString("('a'),")
	}
	values := buffer.String()
	return values[:len(values)-1]
}

func main() {
	port := flag.Int("port", 8080, "The port to listen")
	flag.Parse()
	storeTypes()

	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
