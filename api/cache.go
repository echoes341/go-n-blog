// https://gist.github.com/ismasan/d03d602b8e4e37862547e9a6f0391dc9
// https://gist.github.com/alxshelepenok/0d5c2fb110e19203655e04f4a52e9d87
package main

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

type cachedResponseWriter struct {
	buff  io.Writer
	r     http.ResponseWriter
	multi io.Writer
}

func newCachedResponseWriter(buff io.Writer, resp http.ResponseWriter) http.ResponseWriter {
	multi := io.MultiWriter(buff, resp)
	return &cachedResponseWriter{
		buff:  buff,
		r:     resp,
		multi: multi,
	}
}

func (w *cachedResponseWriter) Header() http.Header {
	return w.r.Header()
}

func (w *cachedResponseWriter) Write(b []byte) (int, error) {
	// here I can intercept body
	return w.multi.Write(b)
}

func (w *cachedResponseWriter) WriteHeader(i int) {
	// here i can intercept header
	w.r.WriteHeader(i)
}

func newCache() {
	cHead = c.New(expiration, 2*expiration)
	cBody = c.New(expiration, 2*expiration)
}

func cacheMdl(fn http.HandlerFunc) http.HandlerFunc {
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
			cw := newCachedResponseWriter(body, w)
			fn(cw, r)

			h := cw.Header() // pick the header, even if it's already sent
			cHead.Set(u, h, c.DefaultExpiration)
			cBody.Set(u, body.String(), c.DefaultExpiration)
		}
	}

}

func removeCacheArticle(url string) {
	cHead.Delete(url)
	cBody.Delete(url)
	log.Printf("[CACHE] Forced cache reloading of %s\n", url)
}
