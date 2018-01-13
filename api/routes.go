package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

func defineRoutes(router *httprouter.Router) {
	router.GET("/article/:id", gzipMdl(cacheMdl(fetchArt)))
	router.GET("/article/:id/likes", gzipMdl(cacheMdl(fetchArtLikes)))
	router.GET("/article/:id/comments", gzipMdl(cacheMdl(fetchArtComments)))
	router.GET("/articles/count", gzipMdl(cacheMdl(countArticles)))
	router.GET("/articles/list", gzipMdl(fetchArticleList))
	router.GET("/articles/list/:year/:month/:day", fetchArticleList)
	router.GET("/test/cache/date", cacheMdl(dateTest))
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

	// prototype
	answer := []map[string]interface{}{}
	/*
		single := map[string]interface{}{}
		single["likes"] = 2
		single["comments"] = 3
		single["article"] = Article{}
		answer = append(answer, single)
		single = map[string]interface{}{} // reset map!
		single["likes"] = 4
		single["comments"] = 5
		single["article"] = Article{}
	*/
	// get parameters handling

	n := 5
	nPar, err := strconv.Atoi(r.FormValue("n"))
	if err == nil && n > 0 && n < 10 { //if n is too great, fall to default
		n = nPar
	}

	likes := r.FormValue("likes") == "true"
	comments := r.FormValue("comments") == "true"
	log.Printf("n: %d l: %v c: %v", n, likes, comments)

	date := time.Now()
	xa := getArticles(n, date)
	for _, article := range xa {
		single := map[string]interface{}{}
		single["article"] = article
		answer = append(answer, single)
	}
	sendJSON(answer, http.StatusOK, w, r)
}
