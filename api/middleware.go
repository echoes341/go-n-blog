package main

import (
	"net/http"
)

// useGET: generic middleware handler for GET methods
type getInterface interface {
	GET(string, http.HandlerFunc)
}

type middleFunc func(http.HandlerFunc) http.HandlerFunc

func useGET(r getInterface, fn middleFunc) *middleWare {
	return &middleWare{fn, r}
}

type middleWare struct {
	fn middleFunc
	getInterface
}

func (m *middleWare) GET(path string, fn http.HandlerFunc) {
	m.getInterface.GET(path, m.fn(fn))
}
