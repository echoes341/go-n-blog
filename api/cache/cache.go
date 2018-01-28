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
	expiration = 5 * time.Minute
)

var cHead *c.Cache
var cBody *c.Cache

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
	cHead = c.New(expiration, 2*expiration)
	cBody = c.New(expiration, 2*expiration)
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
		if head, hFound := cHead.Get(u); hFound {
			log.Printf("[CACHE] Cache present for %s\n", u)
			h := head.(http.Header)
			for key, value := range h {
				w.Header().Set(key, strings.Join(value, ","))
			}

			if body, bFound := cBody.Get(u); bFound {
				fmt.Fprintf(w, "%s", body.(string))
			}
		} else {
			// intercept function writers and save it in the cache
			log.Printf("[CACHE] Building cache for %s\n", u)

			body := bytes.NewBuffer([]byte{})
			cw := newResponseWriter(body, w)
			fn(cw, r)

			h := cw.Header() // pick the header, even if it's already sent
			cHead.Set(u, h, c.DefaultExpiration)
			cBody.Set(u, body.String(), c.DefaultExpiration)
		}
	}

}

// RemoveURL removes given path from cache
func RemoveURL(url string) {
	cHead.Delete(url)
	cBody.Delete(url)
	log.Printf("[CACHE] Forced cache reloading of %s\n", url)
}
