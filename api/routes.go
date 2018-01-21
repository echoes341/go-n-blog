package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dimfeld/httptreemux"
)

func defineRoutes(router *httptreemux.ContextMux) {

	// Single article group -- gzip middleware
	a := useGET(router.NewGroup("/article"), gzipMdl)
	{
		a.GET("/:id", cacheMdl(fetchArt))
		a.GET("/:id/likes", cacheMdl(fetchArtLikes))
		a.GET("/:id/comments", cacheMdl(fetchArtComments))
	}

	// Multiple articles group -- gzip middleware
	xa := useGET(router.NewGroup("/articles"), gzipMdl)
	{
		xa.GET("/count", cacheMdl(countArticles))
		xa.GET("/list", fetchArticleList)
		xa.GET("/list/:year", fetchArticleList)
		xa.GET("/list/:year/:month", fetchArticleList)
		xa.GET("/list/:year/:month/:day", fetchArticleList)
	}

	router.GET("/test/cache/date", cacheMdl(dateTest))

	router.GET("/login", login)
}

func fetchArt(w http.ResponseWriter, r *http.Request) {
	p := httptreemux.ContextParams(r.Context())
	id, err := strconv.Atoi(p["id"])

	if err != nil {
		sendJSON("ID not valid", http.StatusNotFound, w)
		return
	}

	article, err := getArticle(id)
	if err != nil {
		log.Println(err)
		sendJSON("Article not found", http.StatusNotFound, w)
		return
	}

	sendJSON(article, http.StatusOK, w)
}

func fetchArtComments(w http.ResponseWriter, r *http.Request) {
	p := httptreemux.ContextParams(r.Context())
	IDArt, err := strconv.Atoi(p["id"])
	if err != nil {
		sendJSON("ID not valid", http.StatusNotFound, w)
		return
	}

	comments, err := getComments(IDArt)
	if err != nil {
		log.Println(err)
		sendJSON("Comments not found", http.StatusInternalServerError, w)
		return
	}

	if len(comments) == 0 {

		sendJSON("Comments not found", http.StatusNotFound, w)
		return
	}

	sendJSON(comments, http.StatusOK, w)
}

func fetchArtLikes(w http.ResponseWriter, r *http.Request) {
	p := httptreemux.ContextParams(r.Context())
	IDArt, err := strconv.Atoi(p["id"])
	if err != nil {
		sendJSON("ID not valid", http.StatusNotFound, w)
		return
	}

	likes, err := getLikes(IDArt)
	if err != nil {
		log.Println(err)
		sendJSON("Likes not found", http.StatusInternalServerError, w)
		return
	}

	if len(likes) == 0 {
		sendJSON("Likes not found", http.StatusNotFound, w)
		return
	}

	sendJSON(likes, http.StatusOK, w)

}

func dateTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Time now is %d", time.Now().UnixNano())
}

func countArticles(w http.ResponseWriter, r *http.Request) {
	count := getArticleCountByYM()
	sendJSON(count, http.StatusOK, w)
}

func fetchArticleList(w http.ResponseWriter, r *http.Request) {
	// by get params:
	// - n: number of articles to fetch (max 10), default 5
	// - likes: {true/false} count likes for article
	// - comments: {true/false} count comments for article

	// by htttreemux params
	// year, month, day

	p := httptreemux.ContextParams(r.Context())
	answer := []map[string]interface{}{}
	var date time.Time

	year, err := strconv.Atoi(p["year"])
	if err != nil {
		// year it's unreadable
		date = time.Now()
	} else {
		month, err := strconv.Atoi(p["month"])
		if err != nil {
			// month ureadable, default
			month = 12 // December
		} else {
			// dates are in input in js format
			// so month starts from 0
			// but go internal dates start from 1
			// so we need to increment the input
			month++
		}
		day, err := strconv.Atoi(p["day"])
		if err != nil {
			// day unreadable, default
			// from godoc.org/time
			// The month, day, hour, min, sec, and nsec values may be outside their usual ranges
			// and will be normalized during the conversion. For example, October 32 converts to November 1.
			month++
			day = 1
		}
		// set date with given inputsex
		date = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	}

	// get parameters handling
	n := 5
	nPar, err := strconv.Atoi(r.FormValue("n"))
	if err == nil && n > 0 && n < 10 { //if n is too great, fall to default
		n = nPar
	}

	likes := r.FormValue("likes") == "true"
	comments := r.FormValue("comments") == "true"

	xa := getArticles(n, date)
	for _, ar := range xa {
		single := map[string]interface{}{}
		single["article"] = ar
		if likes { // get likes count
			single["likes"] = getLikesCount(ar.ID)
		}
		if comments {
			single["comments"] = getCommentCount(ar.ID)
		}
		answer = append(answer, single)
	}
	if len(answer) == 0 {
		sendJSON(answer, http.StatusNotFound, w)
	} else {
		sendJSON(answer, http.StatusOK, w)
	}
}

func login(w http.ResponseWriter, r *http.Request) {

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
			if verification, ok := err.(jwt.ValidationError); ok { //check if we can see the error as a validation one
				if verification.Errors == jwt.ValidationErrorExpired {
					unauthorized(jwtExpired, w)
					return
				}
			}
			unauthorized(jwtNotValid, w)
			return
		}
		sendJSON("Welcome "+u.Username, http.StatusOK, w)
	}

}
