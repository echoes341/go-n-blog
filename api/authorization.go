package main

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

const (
	badReq byte = iota
	notAuth
)

// https://www.owasp.org/index.php/REST_Security_Cheat_Sheet
// http://blog.restcase.com/restful-api-authentication-basics/

type userToken struct {
	User
	jwt.StandardClaims
}

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

func buildJWT(u User) (string, error) {
	claims := userToken{
		u,
		jwt.StandardClaims{
			ExpiresAt: 150000,
			Issuer:    "gonblog",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte("test"))
	return ss, err
}
