package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/echoes341/go-n-blog/api/models"
)

// LikeController is the basic struct to use all the like-related methods
type LikeController struct{}

// NewLikeController returns an empty LikeController
func NewLikeController() *LikeController {
	return &LikeController{}
}

// Likes is the http handler to get the likes associated with an article
func Likes(w http.ResponseWriter, r *http.Request) {
	p := httptreemux.ContextParams(r.Context())
	IDArt, err := strconv.Atoi(p["id"])
	if err != nil {
		sendJSON("ID not valid", http.StatusNotFound, w)
		return
	}

	l, err := models.Likes(IDArt)
	if err != nil {
		log.Println(err)
		sendJSON("Likes not found", http.StatusInternalServerError, w)
		return
	}

	if len(l) == 0 {
		sendJSON("Likes not found", http.StatusNotFound, w)
		return
	}

	sendJSON(l, http.StatusOK, w)

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
