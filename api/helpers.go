package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type answer struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// Send data as standard JSON (answer struct)
func sendJSON(msg interface{}, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	a := answer{
		Status: status,
		Data:   msg,
	}
	uj, _ := json.Marshal(a)
	fmt.Fprintf(w, "%s", uj)
}

// GzipResponseWriter gzip compress as in https://github.com/socialradar/go-gzip-middleware/blob/master/gzip.go
type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
