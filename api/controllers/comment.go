package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/dimfeld/httptreemux"
	"github.com/echoes341/go-n-blog/api/models"
)

// CommentController is the struct to use all the comments-related methods
type CommentController struct{}

// NewCommentController returns an empty CommentController
func NewCommentController() *CommentController {
	return new(CommentController)
}

// ByArticleID handles request and gives back comments list of an Article
func (uc *CommentController) ByArticleID(w http.ResponseWriter, r *http.Request) {
	p := httptreemux.ContextParams(r.Context())
	IDArt, err := strconv.Atoi(p["id"])
	if err != nil {
		sendJSON("ID not valid", http.StatusNotFound, w)
		return
	}

	c, err := models.Comments(IDArt)
	if err != nil {
		log.Println(err)
		sendJSON("Comments not found", http.StatusInternalServerError, w)
		return
	}

	if len(c) == 0 {
		sendJSON("Comments not found", http.StatusNotFound, w)
		return
	}

	sendJSON(c, http.StatusOK, w)
}
