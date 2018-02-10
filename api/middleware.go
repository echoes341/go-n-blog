package main

import (
	"net/http"
)

type getter interface {
	GET(string, http.HandlerFunc)
}

type poster interface {
	POST(string, http.HandlerFunc)
}

type putter interface {
	PUT(string, http.HandlerFunc)
}

type deleter interface {
	DELETE(string, http.HandlerFunc)
}

type handler interface {
	getter
	poster
	putter
	deleter
}

type middleFunc func(http.HandlerFunc) http.HandlerFunc

// useMdl: generic middleware handler for implemented methods
func useMdl(r handler, fn middleFunc) *middleWare {
	return &middleWare{r, fn}
}

type middleWare struct {
	handler
	fn middleFunc
}

func (m *middleWare) GET(path string, fn http.HandlerFunc) {
	m.handler.GET(path, m.fn(fn))
}
func (m *middleWare) POST(path string, fn http.HandlerFunc) {
	m.handler.POST(path, m.fn(fn))
}
func (m *middleWare) PUT(path string, fn http.HandlerFunc) {
	m.handler.PUT(path, m.fn(fn))
}
func (m *middleWare) DELETE(path string, fn http.HandlerFunc) {
	m.handler.DELETE(path, m.fn(fn))
}
