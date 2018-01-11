package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	c "github.com/patrickmn/go-cache"
)

const (
	expiration = 5 * time.Minute
)

var cHead *c.Cache
var cBody *c.Cache

func newCache() {
	cache = c.New(expiration, 2*expiration)
}

func cacheMdl(fn httprouter.Handle) httprouter.Handle {
	// check if the url is in the cache
	// if yes: call the cache
	// if not: execute function
	//         store the value on the cache
	//         write it on the response
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		u := r.URL.EscapedPath()
		// check cache
		if head, hFound := cHead.Get(u); hFound {
			fmt.Fprintf(w.Header(), "%s", head)
			if body, bFound := cBody.Get(u); bFound {
				fmt.Fprint(w, "%s", body)
			}
		} else {
			// intercept function writers and save it in the cache
		}
	}

}
