package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

func render(w http.ResponseWriter, r *http.Request, statusCode int, content interface{}) {
	var res []byte
	format, contentType := getFormat(r)
	switch format {
	case "json":
		res, _ = json.Marshal(content)
	case "xml":
		res, _ = xml.MarshalIndent(content, "", "  ")
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
