// Package cache contains all the cache-related methods and utilities
// https://gist.github.com/ismasan/d03d602b8e4e37862547e9a6f0391dc9
// https://gist.github.com/alxshelepenok/0d5c2fb110e19203655e04f4a52e9d87
package cache

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	c "github.com/patrickmn/go-cache"
)

const (
	exp = 5 * time.Minute // cache expiration
)

var head *c.Cache
var body *c.Cache

type responseWriter struct {
	buff  io.Writer
	r     http.ResponseWriter
	multi io.Writer
}

func newResponseWriter(buff io.Writer, resp http.ResponseWriter) http.ResponseWriter {
	multi := io.MultiWriter(buff, resp)
	return &responseWriter{
		buff:  buff,
		r:     resp,
		multi: multi,
	}
}

func (w *responseWriter) Header() http.Header {
	return w.r.Header()
}

func (w *responseWriter) Write(b []byte) (int, error) {
	// here I can intercept body
	return w.multi.Write(b)
}

func (w *responseWriter) WriteHeader(i int) {
	// here i can intercept header
	w.r.WriteHeader(i)
}

// Start initialises cache system
func Start() {
	head, body = c.New(exp, 2*exp), c.New(exp, 2*exp)
}

// Middleware it's a http middleware, recording the output of a HandlerFunc
// and repeating it until it's in the cache
func Middleware(fn http.HandlerFunc) http.HandlerFunc {
	// check if the url is in the cache
	// if yes: call the cache
	// if not: execute function
	//         store the value on the cache
	//         write it on the response
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.URL.EscapedPath()
		// check cache
		if hdr, hOk := head.Get(u); hOk {
			log.Printf("[CACHE] Cache present for %s\n", u)
			h := hdr.(http.Header)
			for key, value := range h {
				w.Header().Set(key, strings.Join(value, ","))
			}

			if b, bOk := body.Get(u); bOk {
				fmt.Fprintf(w, "%s", b.(string))
			}
		} else {
			// intercept function writers and save it in the cache
			log.Printf("[CACHE] Building cache for %s\n", u)

			b := bytes.NewBuffer([]byte{})
			cw := newResponseWriter(b, w)
			fn(cw, r)

			h := cw.Header() // pick the header, even if it's already sent
			head.Set(u, h, c.DefaultExpiration)
			body.Set(u, b.String(), c.DefaultExpiration)
		}
	}

}

// RemoveURL removes given path from cache
func RemoveURL(url string) {
	head.Delete(url)
	body.Delete(url)
	log.Printf("[CACHE] Forced cache reloading of %s\n", url)
}
