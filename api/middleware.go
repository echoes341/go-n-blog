package main

import (
	"net/http"
)

// Getter interface containing the GET method
type Getter interface {
	GET(string, http.HandlerFunc)
}

type middleFunc func(http.HandlerFunc) http.HandlerFunc

// useGET: generic middleware handler for GET methods
func useGET(r Getter, fn middleFunc) *middleWare {
	return &middleWare{r, fn}
}

type middleWare struct {
	Getter
	fn middleFunc
}

func (m *middleWare) GET(path string, fn http.HandlerFunc) {
	m.Getter.GET(path, m.fn(fn))
}
