package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/gzip"
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
	store := persistence.NewInMemoryStore(3 * time.Minute)
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(cache.SiteCache(store, 3*time.Minute))
	defineRoutes(router)
	router.Run(":8081")
}
