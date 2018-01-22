package main

import (
	"compress/gzip"
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
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

// gzip middleware
func gzipMdl(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := GzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

// authorization middleware
func authRequired(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// https://gist.github.com/elithrar/9146306
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication
		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 {
			unauthorized(badReq, w)
			return
		}

		switch auth[0] {
		case "Basic":
			b64, err := base64.StdEncoding.DecodeString(auth[1])
			if err != nil {
				unauthorized(badReq, w)
				return
			}

			authDatas := strings.SplitN(string(b64), ":", 2)
			if len(authDatas) != 2 {
				unauthorized(badReq, w)
				return
			}

			user := authDatas[0]
			password := authDatas[1]
			if user == "" || password == "" {
				unauthorized(badReq, w)
				return
			}
			/* SIGNUP
			hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("[FATAL] bcrypt login: %s", err)
				sendJSON(nil, http.StatusInternalServerError, w)
				return
			}*/

			u, err := match(user, password)
			if err != nil {
				unauthorized(userPassNotValid, w)
				return
			}
			token, err := buildJWT(u)
			if err != nil {
				sendJSON(nil, http.StatusInternalServerError, w)
				return
			}
			sendJSON(token, http.StatusOK, w)
		case "Bearer":
			// second argument is JWT
			u, err := checkJWT(auth[1])
			if err != nil {
				log.Printf("%v\n", err)
				// if strings.Contains(err.Error(), "expired") {
				if validation, ok := err.(*jwt.ValidationError); ok {
					if validation.Errors&jwt.ValidationErrorExpired != 0 {
						unauthorized(jwtExpired, w)
						return
					}
				}
				unauthorized(jwtNotValid, w)
				return
			}
			// Put user into context
			r = r.WithContext(addUserToContext(r.Context(), &u))
			fn(w, r)
		default:
			unauthorized(badReq, w)
		}

	}
}
