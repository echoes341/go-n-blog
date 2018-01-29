package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/echoes341/go-n-blog/api/cache"
	"github.com/echoes341/go-n-blog/api/models"
)

// ArticleController is the specific structure controlling articles
type ArticleController struct{}

// NewArticleController returns an emptty ArticleController to use the specific article methods
func NewArticleController() *ArticleController {
	return &ArticleController{}
}

// Fetch is a http request handler to returns an article as JSON
func (ac *ArticleController) Fetch(w http.ResponseWriter, r *http.Request) {
	p := httptreemux.ContextParams(r.Context())
	id, err := strconv.Atoi(p["id"])

	if err != nil {
		sendJSON("ID not valid", http.StatusNotFound, w)
		return
	}

	article, err := models.ArticleGet(id)
	if err != nil {
		log.Println(err)
		sendJSON("Article not found", http.StatusNotFound, w)
		return
	}

	sendJSON(article, http.StatusOK, w)
}

// Count is a http handler which returns the article count as descripted in models as a JSON
func (ac *ArticleController) Count(w http.ResponseWriter, r *http.Request) {
	count := models.ArticleCount()
	sendJSON(count, http.StatusOK, w)
}

// List gives n articles before the date selected
// Default values:
// n: 5
// date: date.Now()
func (ac *ArticleController) List(w http.ResponseWriter, r *http.Request) {
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
	if err == nil && n > 0 && n < 10 { //if n is too great or not valid, fall to default
		n = nPar
	}

	likes := r.FormValue("likes") == "true"
	comments := r.FormValue("comments") == "true"

	xa := models.Articles(n, date)
	for _, ar := range xa {
		single := map[string]interface{}{}
		single["article"] = ar
		if likes { // get likes count
			single["likes"] = models.LikesCount(ar.ID)
		}
		if comments {
			single["comments"] = models.CommentsCount(ar.ID)
		}
		answer = append(answer, single)
	}
	if len(answer) == 0 {
		sendJSON(answer, http.StatusNotFound, w)
	} else {
		sendJSON(answer, http.StatusOK, w)
	}
}

// Edit edits an article, should be used by the authorization package
func (ac *ArticleController) Edit(w http.ResponseWriter, r *http.Request) {
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

		a := models.ArticleUpdate(id, title, text)
		if a.ID == 0 { // Something went wrong
			sendJSON("Error: impossible to edit article", http.StatusInternalServerError, w)
			return
		}

		url := r.URL.EscapedPath()
		cache.RemoveURL(url)

		w.Header().Set("Content-Location", url)
		sendJSON(a, http.StatusOK, w)
	} else {
		sendJSON("You are not admin", http.StatusForbidden, w)
	}
}

// Add is the http handler to add an Article
func (ac *ArticleController) Add(w http.ResponseWriter, r *http.Request) {
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

		a := models.ArticleAdd(title, text, author, date)
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

// Delete is the http handler to remove an Article
func (ac *ArticleController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := userContext(ctx)
	if !u.IsAdmin {
		sendJSON("You are not admin", http.StatusForbidden, w)
		return
	}

	p := httptreemux.ContextParams(ctx)
	idParam, _ := strconv.Atoi(p["id"])
	if idParam <= 0 {
		sendJSON("Input not valid", http.StatusBadRequest, w)
		return
	}

	id := uint(idParam)
	notFound, err := models.ArticleRemove(id)
	if err != nil {
		if notFound {
			sendJSON("Article not found", http.StatusNotFound, w)
			return
		}
		sendJSON("Internal Error", http.StatusInternalServerError, w)
		return
	}
	cache.RemoveURL(r.URL.EscapedPath())

	sendJSON("Delete ok", http.StatusOK, w)
}
