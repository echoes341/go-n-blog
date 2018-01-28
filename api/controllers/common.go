package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
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
