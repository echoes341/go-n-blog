package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/echoes341/go-n-blog/api/cache"
	"github.com/echoes341/go-n-blog/api/controllers"
)

const (
	articleGroup  = "/article"
	articlesGroup = "/articles"
)

func defineRoutes(router *httptreemux.ContextMux) {

	// Single article group -- gzip middleware
	a := router.NewGroup(articleGroup)
	agz := useGET(a, gzipMdl)
	ac := controllers.NewArticleController()
	{
		agz.GET("/:id", cache.Middleware(ac.Fetch))
		agz.GET("/:id/likes", cache.Middleware(fetchArtLikes))
		agz.GET("/:id/comments", cache.Middleware(fetchArtComments))
	}

	// Reserved section

	{
		// add article
		a.POST("/", ac.Add)
		// edit article
		a.PUT("/:id", ac.Edit)
		// remove article
		a.DELETE("/:id", ac.Delete)
	}

	// Multiple articles group -- gzip middleware
	xa := useGET(router.NewGroup(articlesGroup), gzipMdl)
	{
		// count article
		xa.GET("/count", cache.Middleware(ac.Count))
		// get articles by date
		xa.GET("/list", ac.List)
		xa.GET("/list/:year", ac.List)
		xa.GET("/list/:year/:month", ac.List)
		xa.GET("/list/:year/:month/:day", ac.List)
	}

	router.GET("/test/cache/date", cache.Middleware(dateTest))

	router.GET("/login", authRequired(login))
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

func login(w http.ResponseWriter, r *http.Request) {
	u := userContext(r.Context())
	var msg string
	if u.IsAdmin {
		msg = fmt.Sprintf("Welcome %s. You are admin!", u.Username)
	} else {
		msg = fmt.Sprintf("Welcome %s. You are not admin", u.Username)
	}
	sendJSON(msg, http.StatusOK, w)
}
