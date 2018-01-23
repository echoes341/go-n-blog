package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dimfeld/httptreemux"
)

const (
	articleGroup  = "/article"
	articlesGroup = "/articles"
)

func defineRoutes(router *httptreemux.ContextMux) {

	// Single article group -- gzip middleware
	a := router.NewGroup(articleGroup)
	agz := useGET(a, gzipMdl)
	{
		agz.GET("/:id", cacheMdl(fetchArt))
		agz.GET("/:id/likes", cacheMdl(fetchArtLikes))
		agz.GET("/:id/comments", cacheMdl(fetchArtComments))
	}

	// Reserved section

	{
		a.POST("/", addArticleRoute)
		a.PUT("/:id", editArticle)
	}

	// Multiple articles group -- gzip middleware
	xa := useGET(router.NewGroup(articlesGroup), gzipMdl)
	{
		xa.GET("/count", cacheMdl(countArticles))
		xa.GET("/list", fetchArticleList)
		xa.GET("/list/:year", fetchArticleList)
		xa.GET("/list/:year/:month", fetchArticleList)
		xa.GET("/list/:year/:month/:day", fetchArticleList)
	}

	router.GET("/test/cache/date", cacheMdl(dateTest))

	router.GET("/login", authRequired(login))
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
	u := userContext(r.Context())
	var msg string
	if u.IsAdmin {
		msg = fmt.Sprintf("Welcome %s. You are admin!", u.Username)
	} else {
		msg = fmt.Sprintf("Welcome %s. You are not admin", u.Username)
	}
	sendJSON(msg, http.StatusOK, w)
}

func addArticleRoute(w http.ResponseWriter, r *http.Request) {
	u := userContext(r.Context())

	if u.IsAdmin {
		title := r.FormValue("title")
		aID, _ := strconv.Atoi(r.FormValue("author"))

		text := r.FormValue("text")
		dateInt, _ := strconv.Atoi(r.FormValue("date")) // date in unix format

		if title == "" || text == "" {
			sendJSON("Input not valid", http.StatusBadRequest, w)
			return
		}

		var author uint
		if aID <= 0 {
			author = u.ID
		} else {
			author = uint(aID)
		}
		var date time.Time
		if dateInt == 0 { // Error in conversion or parameter empty
			date = time.Now()
		} else {
			date = time.Unix(int64(dateInt), 0)
		}

		a := addArticle(title, text, author, date)
		if a.ID == 0 { // Something went wrong
			sendJSON("Error: impossible to add article", http.StatusInternalServerError, w)
			return
		}

		w.Header().Set("Content-Location", fmt.Sprintf("%s/%d", articleGroup, a.ID))
		sendJSON(a, http.StatusCreated, w)
	} else {
		sendJSON("You are not admin", http.StatusForbidden, w)
	}
}

func editArticle(w http.ResponseWriter, r *http.Request) {
	// u := userContext(r.Context())
	// debug: dummy user
	ctx := r.Context()
	u := userContext(ctx)
	if u.IsAdmin {
		p := httptreemux.ContextParams(ctx)

		idParam, _ := strconv.Atoi(p["id"])
		if idParam <= 0 { // conversion failed or bad input
			sendJSON("Input not valid", http.StatusBadRequest, w)
			return
		}

		id := uint(idParam)
		title := r.FormValue("title")
		text := r.FormValue("text")

		a := updateArticle(id, title, text)
		if a.ID == 0 { // Something went wrong
			sendJSON("Error: impossible to edit article", http.StatusInternalServerError, w)
			return
		}

		url := fmt.Sprintf("%s/%d", articleGroup, a.ID)

		removeCacheArticle(url)

		w.Header().Set("Content-Location", url)
		sendJSON(a, http.StatusOK, w)
	} else {
		sendJSON("You are not admin", http.StatusForbidden, w)
	}
}

