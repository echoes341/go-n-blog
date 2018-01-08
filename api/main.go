package main

import (
	"log"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	// open db connection
	var err error
	db, err = gorm.Open("mysql", "gonblog:gonblog@tcp(127.0.0.1:3306)/gonblog?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Panicln("Database connection failed")
	}

	db.AutoMigrate(&articleDB{})
	db.AutoMigrate(&commentDB{})
	db.AutoMigrate(&likeDB{})
}

func main() {
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	defineRoutes(router)
	router.Run(":8081")
}
