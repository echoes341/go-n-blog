package main

import (
	"log"
	"net/http"
)

const (
	badReq byte = iota
	notAuth
)

func unauthorized(status byte, w http.ResponseWriter) {
	var message string
	switch status {
	case badReq:
		message = "Bad Request"
	case notAuth:
		message = "Username and/or password do not match"
	}
	log.Printf("LOGIN Forbidden: %s", message)
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	sendJSON(message, http.StatusUnauthorized, w)
}
