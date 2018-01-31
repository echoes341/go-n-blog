package controllers

import (
	"log"
	"net/http"
	"strconv"

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
func (lc *LikeController) Likes(w http.ResponseWriter, r *http.Request) {
	p := httptreemux.ContextParams(r.Context())
	IDArt, err := strconv.Atoi(p["id"])
	if err != nil {
		sendJSON(ErrIDNotValid, http.StatusNotFound, w)
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

// Toggle is http handler to toggle like on given article
func (lc *LikeController) Toggle(w http.ResponseWriter, r *http.Request) {
	p := httptreemux.ContextParams(r.Context())
	aID, _ := strconv.Atoi(p["id"])

	if aID <= 0 {
		sendJSON(ErrIDNotValid, http.StatusBadRequest, w)
		return
	}
	u := models.UserContext(r.Context())
	added, err := models.LikeToggle(uint(aID), u.ID)
	if err != nil {
		sendJSON(err, http.StatusInternalServerError, w)
		return
	}
	sendJSON(added, http.StatusOK, w)
}
