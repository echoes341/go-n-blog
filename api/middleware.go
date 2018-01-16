package main

import (
	"net/http"
)

// GET middleware: gzip enabled

func gzEnable(r getInterface) *gzGroup {
	return &gzGroup{r}
}

type gzGroup struct {
	getInterface
}

type getInterface interface {
	GET(string, http.HandlerFunc)
}

func (g *gzGroup) GET(path string, fn http.HandlerFunc) {
	g.getInterface.GET(path, gzipMdl(fn))
}
