package controllers

import (
	"net/http"
	"strconv"

	"github.com/dimfeld/httptreemux"

	"github.com/echoes341/go-n-blog/api/models"
)

// UserController is struct to control all methods associate with user
type UserController struct{}

// NewUserController returns a new UserController
func NewUserController() *UserController {
	return new(UserController)
}

// SignUp is the handler to sign a new user in the system
func (uc *UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	// taking infos from post parameters
	u := r.FormValue("username")
	m := r.FormValue("email")
	p := r.FormValue("password")

	if m == "" || u == "" || p == "" {
		sendJSON(ErrLoginBadRequest.Error(), http.StatusBadRequest, w)
		return
	}

	User, err := models.UserAdd(u, m, p)
	if err != nil {
		var status int
		if err == models.ErrUserPresent {
			status = http.StatusConflict
		} else {
			status = http.StatusInternalServerError
		}
		sendJSON(err.Error(), status, w)
		return
	}
	sendJSON(User, http.StatusCreated, w)
}

// Remove is a function to remove articles
func (uc *UserController) Remove(w http.ResponseWriter, r *http.Request) {
	par := httptreemux.ContextParams(r.Context())
	id, _ := strconv.Atoi(par["id"])
	if id <= 0 {
		sendJSON(ErrIDNotValid.Error(), http.StatusBadRequest, w)
		return
	}

	err := models.UserRemove(uint(id))
	if err != nil {
		sendJSON(err.Error(), http.StatusInternalServerError, w)
		return
	}

	sendJSON(nil, http.StatusNoContent, w)
}
