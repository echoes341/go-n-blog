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
	likeGroup     = "/like"
	userGroup     = "/user"
)

func defineRoutes(router *httptreemux.ContextMux) {

	ac := controllers.NewArticleController()
	lc := controllers.NewLikeController()
	cc := controllers.NewCommentController()
	at := controllers.NewAuth()
	uc := controllers.NewUserController()

	a := router.NewGroup(articleGroup)
	// Single article group -- gzip middleware
	aGz := useMdl(
		useMdl(a, CORSAll),
		gzipMdl,
	)
	{
		// get article by id
		aGz.GET("/:id", cache.Middleware(ac.Fetch))
		// get related likes of an article
		aGz.GET("/:id/likes", cache.Middleware(lc.Likes))
		// get related comments of an article
		aGz.GET("/:id/comments", cache.Middleware(cc.ByArticleID))
	}

	l := router.NewGroup(likeGroup)
	{
		l.POST("/:id", lc.Toggle)
	}

	// Reserved section
	{
		// add article
		a.POST("", at.ExecIfAdmin(ac.Add))
		// edit article
		a.PUT("/:id", at.ExecIfAdmin(ac.Edit))
		// remove article
		a.DELETE("/:id", at.ExecIfAdmin(ac.Delete))
		// post a comment
		a.POST("/:id/comment", at.AuthRequired(cc.Add))
	}

	// Multiple articles group -- gzip middleware
	asGz := useMdl(router.NewGroup(articlesGroup), gzipMdl)
	{
		// count article
		asGz.GET("/count", cache.Middleware(ac.Count))
		// get articles by date
		asGz.GET("/list", ac.List)
		asGz.GET("/list/:year", ac.List)
		asGz.GET("/list/:year/:month", ac.List)
		asGz.GET("/list/:year/:month/:day", ac.List)
	}

	// Debug routes
	router.GET("/test/cache/date", cache.Middleware(dateTest))
	router.GET("/login", at.AuthRequired(login))
	// User management
	ur := router.NewGroup(userGroup)
	{
		// signin
		ur.POST("", uc.SignUp)
		ur.DELETE("/:id", at.ExecIfAdmin(uc.Remove))
	}
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
