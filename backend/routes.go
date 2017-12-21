package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func initRoutes(m *httprouter.Router) {
	m.GET("/", home)
	m.GET("/insert", userTestInsert) // !!! DOES NOT WORK
	m.GET("/get/users", printUsers)
	m.GET("/login", login)
	m.ServeFiles("/css/*filepath", http.Dir("./resources/css"))

}
