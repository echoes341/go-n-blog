package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/echoes341/go-n-blog/api/models"
)

var (
	// ErrLoginBadRequest is encountered when login request is wrong
	ErrLoginBadRequest = errors.New("Bad Request")
	// ErrLoginNotValid is encountered when username or password are wrong
	ErrLoginNotValid = models.ErrLoginError
	// ErrJWTNotValid is encountered when JWT examined is malformed
	ErrJWTNotValid = errors.New("JWT not valid")
	// ErrJWTExpired is encountered when JWT is expired
	ErrJWTExpired = errors.New("JWT is expired")
)

// https://www.owasp.org/index.php/REST_Security_Cheat_Sheet
// http://blog.restcase.com/restful-api-authentication-basics/

type userToken struct {
	models.User
	jwt.StandardClaims
}

// Auth is the struct to use the authorization functions
type Auth struct{}

// NewAuth returns a zeroed Auth struct pointer
func NewAuth() *Auth {
	return new(Auth)
}

func unauthorized(err error, w http.ResponseWriter) {
	log.Printf("LOGIN Forbidden: %s", err)
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	sendJSON(err.Error(), http.StatusUnauthorized, w)
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

// ExecIfAdmin checks if login is valid and if the user is admin and then executes the function associated
func (at *Auth) ExecIfAdmin(fn http.HandlerFunc) http.HandlerFunc {
	return at.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
		u := models.UserContext(r.Context())
		if !u.IsAdmin {
			sendJSON("You are not admin", http.StatusForbidden, w)
		}
		fn(w, r)
	})
}

// AuthRequired is an authorization middleware
func (at *Auth) AuthRequired(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// https://gist.github.com/elithrar/9146306
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication
		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 {
			unauthorized(ErrLoginBadRequest, w)
			return
		}

		switch auth[0] {
		// Auth and get jwt
		case "Basic":
			b64, err := base64.StdEncoding.DecodeString(auth[1])
			if err != nil {
				unauthorized(ErrLoginBadRequest, w)
				return
			}

			authDatas := strings.SplitN(string(b64), ":", 2)
			if len(authDatas) != 2 {
				unauthorized(ErrLoginBadRequest, w)
				return
			}

			n := authDatas[0] // username
			p := authDatas[1] // password
			if n == "" || p == "" {
				unauthorized(ErrLoginBadRequest, w)
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
				unauthorized(ErrLoginNotValid, w)
				return
			}
			token, err := buildJWT(u)
			if err != nil {
				sendJSON(nil, http.StatusInternalServerError, w)
				return
			}
			sendJSON(token, http.StatusOK, w)
		// Check jwt and fill user infos
		case "Bearer":
			// second argument is JWT
			u, err := checkJWT(auth[1])
			if err != nil {
				log.Printf("%v\n", err)
				// if strings.Contains(err.Error(), "expired") {
				if validation, ok := err.(*jwt.ValidationError); ok {
					if validation.Errors&jwt.ValidationErrorExpired != 0 {
						unauthorized(ErrJWTExpired, w)
						return
					}
				}
				unauthorized(ErrJWTNotValid, w)
				return
			}
			// Put user into context
			r = r.WithContext(models.UserAddToContext(r.Context(), &u))
			fn(w, r)
		default:
			unauthorized(ErrLoginBadRequest, w)
		}

	}
}
