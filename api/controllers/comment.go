package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"

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

// Add is http handler to add a comment under a specific article
func (uc *CommentController) Add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	IDArt, _ := strconv.Atoi(httptreemux.ContextParams(ctx)["id"])
	if IDArt <= 0 {
		sendJSON(ErrIDNotValid, http.StatusBadRequest, w)
		return
	}
	u := models.UserContext(ctx)
	// take content from post
	// date is set by the server
	ctn := r.FormValue("content")
	if ctn == "" {
		sendJSON(ErrParameterEmpty, http.StatusBadRequest, w)
		return
	}

	d := time.Now()

	c, err := models.CommentAdd(uint(IDArt), u.ID, d, ctn)
	if err != nil {
		var status int
		switch err {
		case gorm.ErrRecordNotFound:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}
		sendJSON(err, status, w)
		return
	}
	// here add is successful
	sendJSON(c, http.StatusOK, w)
}

// Delete is http handler to delete a comment by its own ID
func (uc *CommentController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idC, _ := strconv.Atoi(httptreemux.ContextParams(ctx)["id"])
	if idC <= 0 {
		sendJSON(ErrIDNotValid, http.StatusBadRequest, w)
		return
	}

	err := models.CommentRemove(uint(idC))
	if err != nil {
		sendJSON(err, http.StatusInternalServerError, w)
		return
	}
	sendJSON("DELETE OK", http.StatusOK, w)
}

// Edit is http handler to edit comment by its ID
func (uc *CommentController) Edit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idC, _ := strconv.Atoi(httptreemux.ContextParams(ctx)["id"])
	if idC <= 0 {
		sendJSON(ErrIDNotValid, http.StatusBadRequest, w)
		return
	}

	cnt := r.FormValue("content")

	// every content is accepted?
	// how to manage bot?
	// maybe using recaptcha as middleware?
	// [TODO] an user can edit only his comments
	// except for admin, who can edit every comment he wants

	// u := models.UserContext(ctx)

	c, err := models.CommentUpdate(uint(idC), cnt)
	if err != nil {
		var status int
		switch err {
		case gorm.ErrRecordNotFound:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}
		sendJSON(err, status, w)
		return
	}

	sendJSON(c, http.StatusOK, w)
}
