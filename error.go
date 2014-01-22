package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

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
	errorHandler(
		w,
		r,
		http.StatusBadRequest,
		[]string{err.Error()},
	)
}

func errorHandler(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	errors []string,
) {
	error := ErrorWrapper{status, getStatusMessage(status), errors}

	var res []byte
	format, contentType := getFormat(r)
	switch format {
	case "json":
		res, _ = json.Marshal(error)
	case "xml":
		res, _ = xml.MarshalIndent(error, "", "  ")
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	w.Write(res)
}
