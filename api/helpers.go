package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type answer struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// Send data as standard JSON (answer struct)
func sendJSON(msg interface{}, status int, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	a := answer{
		Status: status,
		Data:   msg,
	}
	uj, _ := json.Marshal(a)
	fmt.Fprintf(w, "%s", uj)
}

// gzip compress as in https://github.com/socialradar/go-gzip-middleware/blob/master/gzip.go
type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipMdl(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r, p)
			return
		}
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := GzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r, p)
	}
}
