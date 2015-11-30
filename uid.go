package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	// imported because is used as a driver
	// read https://golang.org/doc/effective_go.html#blank_unused
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// based on http://code.flickr.net/2010/02/08/ticket-servers-distributed-unique-primary-keys-on-the-cheap/

func handler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			// @todo we should close the db connection
			//if db != nil {
			//	db.Close()
			//}
			log.Println("Panic: ", err)
			errorHandlerFromError(w, r, http.StatusInternalServerError, err.(error))
		}
	}()

	c := conf{
		increment: 10,
		maxAmount: 100,
		tableName: "uid1",
		hosts: map[string]string{
			"first":  "root:root@unix(/opt/local/var/run/mysql56/mysqld.sock)/uid",
			"second": "root:root2@unix(/opt/local/var/run/mysql56/mysqld.sock)/uid",
		},
	}

	// get amount
	var err error
	var amount int
	if amount, err = getAmount(r, c.maxAmount); err != nil {
		errorHandlerFromError(w, r, http.StatusBadRequest, err)
		return
	}

	ids := newIdGenerator(c.getRandomHost(), c.tableName, c.increment).get(amount)

	render(w, r, http.StatusOK, ids)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	errorHandler(
		w,
		r,
		http.StatusNotFound,
		[]string{
			fmt.Sprintf("No route found for \"%s %s\"", r.Method, r.URL.Path[0:]),
		},
	)
}

func getAmount(r *http.Request, maxAmount int) (int, error) {
	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil || (amount < 1 || amount > maxAmount) {
		return 1, errors.New(fmt.Sprintf("Amount \"%s\" is not a valid value", r.FormValue("amount")))
	}

	return amount, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	port := flag.Int("port", 8080, "The port to listen")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/api/ids", handler)
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
}
