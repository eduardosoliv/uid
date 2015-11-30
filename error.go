package main

import "net/http"

type ErrorWrapper struct {
	Status  int
	Message string
	Errors  []string
}

func errorHandlerFromError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	err error,
) {
	errorHandler(w, r, status, []string{err.Error()})
}

func errorHandler(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	errors []string,
) {
	error := ErrorWrapper{status, http.StatusText(status), errors}

	render(w, r, status, error)
}
