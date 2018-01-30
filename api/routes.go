package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/echoes341/go-n-blog/api/cache"
	"github.com/echoes341/go-n-blog/api/controllers"
	"github.com/echoes341/go-n-blog/api/models"
)

const (
	articleGroup  = "/article"
	articlesGroup = "/articles"
)

func defineRoutes(router *httptreemux.ContextMux) {

	// Single article group -- gzip middleware
	a := router.NewGroup(articleGroup)
	ac := controllers.NewArticleController()
	lc := controllers.NewLikeController()
	uc := controllers.NewCommentController()
	at := controllers.NewAuth()
	agz := useGET(a, gzipMdl)
	{
		// get article by id
		agz.GET("/:id", cache.Middleware(ac.Fetch))
		// get related likes of an article
		agz.GET("/:id/likes", cache.Middleware(lc.Likes))
		// get related comments of an article
		agz.GET("/:id/comments", cache.Middleware(uc.ByArticleID))
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

	router.GET("/login", at.AuthRequired(login))
}

func dateTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Time now is %d", time.Now().UnixNano())
}

func login(w http.ResponseWriter, r *http.Request) {
	u := models.UserContext(r.Context())
	var msg string
	if u.IsAdmin {
		msg = fmt.Sprintf("Welcome %s. You are admin!", u.Username)
	} else {
		msg = fmt.Sprintf("Welcome %s. You are not admin", u.Username)
	}
	fmt.Fprintln(w, msg)
}
