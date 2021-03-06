package main

import (
	"log"
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/echoes341/go-n-blog/api/cache"
	"github.com/echoes341/go-n-blog/api/models"
)

func init() {
	// open db connection
	var err error
	err = models.NewDB("gonblog", "gonblog", "127.0.0.1:3306", "gonblog")
	if err != nil {
		log.Panicln("Database connection failed")
	}
}

func main() {
	mux := httptreemux.NewContextMux()
	cache.Start()
	defineRoutes(mux)
	http.ListenAndServe(":8080", mux)
}
