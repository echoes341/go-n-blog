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
