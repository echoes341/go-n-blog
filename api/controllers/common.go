package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrIDNotValid is for not valid article/likes/comments ids ( <= 0 )
	ErrIDNotValid = errors.New("ID not valid")
)

type answer struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func sendJSON(msg interface{}, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	a := answer{
		Status: status,
		Data:   msg,
	}
	uj, _ := json.Marshal(a)
	fmt.Fprintf(w, "%s", uj)
}
