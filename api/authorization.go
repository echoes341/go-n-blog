package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	badReq byte = iota
	userPassNotValid
	jwtNotValid
	jwtExpired
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
	case userPassNotValid:
		message = "Username and/or password do not match"
	case jwtNotValid:
		message = "JWT not valid"
	case jwtExpired:
		message = "JWT expired"
	}
	log.Printf("LOGIN Forbidden: %s", message)
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	sendJSON(message, http.StatusUnauthorized, w)
}

func buildJWT(u User) (string, error) {
	expiration := time.Now().Add(10 * time.Minute)
	claims := userToken{
		u,
		jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
			Issuer:    "gonblog",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte("test"))

	return ss, err
}

func checkJWT(j string) (User, error) {
	token, err := jwt.ParseWithClaims(j, &userToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("test"), nil
	})
	if err != nil {
		return User{}, err
	}

	if claims, ok := token.Claims.(*userToken); ok && token.Valid {
		return claims.User, nil
	}
	return User{}, fmt.Errorf("JWT not valid")
}
