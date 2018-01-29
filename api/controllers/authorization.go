package controllers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/echoes341/go-n-blog/api/models"
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
	models.User
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

func buildJWT(u models.User) (string, error) {
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

func checkJWT(j string) (models.User, error) {
	token, err := jwt.ParseWithClaims(j, &userToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("test"), nil
	})
	if err != nil {
		return models.User{}, err
	}

	if claims, ok := token.Claims.(*userToken); ok && token.Valid {
		return claims.User, nil
	}
	return models.User{}, fmt.Errorf("JWT not valid")
}

func login(w http.ResponseWriter, r *http.Request) {
	u := models.UserContext(r.Context())
	var msg string
	if u.IsAdmin {
		msg = fmt.Sprintf("Welcome %s. You are admin!", u.Username)
	} else {
		msg = fmt.Sprintf("Welcome %s. You are not admin", u.Username)
	}
	sendJSON(msg, http.StatusOK, w)
}

// ExecIfAdmin checks if the user is admin and then executes the function associated
func ExecIfAdmin(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := models.UserContext(r.Context())
		if !u.IsAdmin {
			sendJSON("You are not admin", http.StatusForbidden, w)
		}
		fn(w, r)
	}
}

// AuthRequired is an authorization middleware
func AuthRequired(fn http.HandlerFunc) http.HandlerFunc {
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

			n := authDatas[0] // username
			p := authDatas[1] // password
			if n == "" || p == "" {
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

			u, err := models.UserMatch(n, p)
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
			r = r.WithContext(models.UserAddToContext(r.Context(), &u))
			fn(w, r)
		default:
			unauthorized(badReq, w)
		}

	}
}
