package main

import (
	"log"
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/echoes341/go-n-blog/api/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	// open db connection
	var err error
	db, err = models.NewDB("gonblog", "gonblog", "127.0.0.1:3306", "gonblog")
	if err != nil {
		log.Panicln("Database connection failed")
	}

	db.AutoMigrate(&userDB{})
	db.AutoMigrate(&articleDB{})
	db.AutoMigrate(&commentDB{})
	db.AutoMigrate(&likeDB{})
}

func main() {
	mux := httptreemux.NewContextMux()
	newCache()
	defineRoutes(mux)
	http.ListenAndServe(":8080", mux)
}
