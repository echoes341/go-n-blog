package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/julienschmidt/httprouter"
)

func defineRoutes(router *httprouter.Router) {
	router.GET("/article/:id", gzipMdl(cacheMdl(fetchArt)))
	router.GET("/article/:id/likes", gzipMdl(cacheMdl(fetchArtLikes)))
	router.GET("/article/:id/comments", gzipMdl(cacheMdl(fetchArtComments)))
	router.GET("/articles/count", gzipMdl(cacheMdl(countArticles)))
	router.GET("/articles/list", gzipMdl(fetchArticleList))
	router.GET("/articles/list/:year", gzipMdl(fetchArticleList))
	router.GET("/articles/list/:year/:month", gzipMdl(fetchArticleList))
	router.GET("/articles/list/:year/:month/:day", gzipMdl(fetchArticleList))
	router.GET("/test/cache/date", cacheMdl(dateTest))

	router.GET("/login", login)
}

func fetchArt(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.Atoi(p.ByName("id"))

	if err != nil {
		sendJSON("ID not valid", http.StatusNotFound, w, r)
		return
	}

	article, err := getArticle(id)
	if err != nil {
		log.Println(err)
		sendJSON("Article not found", http.StatusNotFound, w, r)
		return
	}

	sendJSON(article, http.StatusOK, w, r)
}

func fetchArtComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	IDArt, err := strconv.Atoi(p.ByName("id"))
	if err != nil {
		sendJSON("ID not valid", http.StatusNotFound, w, r)
		return
	}

	comments, err := getComments(IDArt)
	if err != nil {
		log.Println(err)
		sendJSON("Comments not found", http.StatusInternalServerError, w, r)
		return
	}

	if len(comments) == 0 {

		sendJSON("Comments not found", http.StatusNotFound, w, r)
		return
	}

	sendJSON(comments, http.StatusOK, w, r)
}

func fetchArtLikes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	IDArt, err := strconv.Atoi(p.ByName("id"))
	if err != nil {
		sendJSON("ID not valid", http.StatusNotFound, w, r)
		return
	}

	likes, err := getLikes(IDArt)
	if err != nil {
		log.Println(err)
		sendJSON("Likes not found", http.StatusInternalServerError, w, r)
		return
	}

	if len(likes) == 0 {
		sendJSON("Likes not found", http.StatusNotFound, w, r)
		return
	}

	sendJSON(likes, http.StatusOK, w, r)

}

func dateTest(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "Time now is %d", time.Now().UnixNano())
}

func countArticles(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	count := getArticleCountByYM()
	sendJSON(count, http.StatusOK, w, r)
}

func fetchArticleList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// by get params:
	// - n: number of articles to fetch (max 10), default 5
	// - likes: {true/false} count likes for article
	// - comments: {true/false} count comments for article

	// by httrouter params
	// year, month, day

	answer := []map[string]interface{}{}
	var date time.Time

	year, err := strconv.Atoi(p.ByName("year"))
	if err != nil {
		// year it's unreadable
		date = time.Now()
	} else {
		month, err := strconv.Atoi(p.ByName("month"))
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
		day, err := strconv.Atoi(p.ByName("day"))
		if err != nil {
			// day readable, default
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
		sendJSON(answer, http.StatusNotFound, w, r)
	} else {
		sendJSON(answer, http.StatusOK, w, r)
	}
}

func login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// GET params user and password
	user := r.FormValue("user")
	password := r.FormValue("password")

	if user == "" || password == "" {
		sendJSON(nil, http.StatusBadRequest, w, r)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[FATAL] bcrypt login: %s", err)
		sendJSON(nil, http.StatusInternalServerError, w, r)
		return
	}

	log.Printf("user: %s\n", user)
	log.Printf("pwd hash: %s", hash)
}
