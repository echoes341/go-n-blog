package main

import (
	"io"
	"net/http"
)

// GzipResponseWriter gzip compress as in https://github.com/socialradar/go-gzip-middleware/blob/master/gzip.go
type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
